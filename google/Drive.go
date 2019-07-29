// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package google

import (
	"github.com/reaction-eng/restlib/configuration"
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"golang.org/x/net/html"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Drive struct {

	//Store the connection to google.  This has been wrapped with the correct Headers
	connection *drive.Service

	//Store the preview length
	previewLength int

	//Store the timezone
	timeZone string
}

//Get a new interface
func NewDrive(configFiles ...string) *Drive {
	//Create a new config
	config, err := configuration.NewConfiguration(configFiles...)

	//Create a new
	gInter := &Drive{
		timeZone: config.GetStringFatal("default_time_zone"),
	}

	//Open the client
	jwtConfig := &jwt.Config{
		Email:      config.GetStringFatal("google_auth_email"),
		PrivateKey: []byte(config.GetStringFatal("google_auth_key")),
		Scopes: []string{
			drive.DriveMetadataReadonlyScope,
			drive.DriveReadonlyScope,
			drive.DriveFileScope,
		},
		TokenURL: google.JWTTokenURL,
	}

	//Build the connection
	httpCon := jwtConfig.Client(context.Background())

	//Now build the drive service
	driveConn, err := drive.New(httpCon)
	gInter.connection = driveConn

	//Check for errors
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	gInter.previewLength, _ = config.GetInt("preview_length")
	return gInter
}

//See if starts with a date and name
func (gog *Drive) splitNameAndDate(nameIn string) (string, *time.Time) {

	//Trime the name in
	nameIn = strings.TrimSpace(nameIn)

	//Take the first part before a space
	splitLoc := strings.Index(nameIn, " ")

	//If there is at least one loc
	if splitLoc >= 0 {
		//Get the first item
		firstPart := nameIn[0:splitLoc]

		//Now keep testing if it is a date
		date, err := time.Parse("2006-1-2", firstPart)
		if err != nil {
			date, err = time.Parse("2006/1/2", firstPart)
		}
		if err != nil {
			date, err = time.Parse("2006.1.2", firstPart)
		}
		if err != nil {
			date, err = time.Parse("2/1/2006", firstPart)
		}

		//If we got a date return it
		if err == nil {
			loc, err := time.LoadLocation(gog.timeZone)

			//convert to the local time zone
			if err == nil {
				//Move the dat
				date = date.In(loc)

				//Get the offset for that time
				_, offset := date.Zone()

				//update the offset
				offset = -offset

				//Add the offset to the date
				date = date.Add(time.Duration(offset) * time.Second)

			}

			return strings.TrimSpace(nameIn[splitLoc:]), &date
		} else {
			//Just return the name
			return nameIn, nil
		}

	} else {
		//Just return the name
		return nameIn, nil
	}

}

/**
Recursive call to build the file list
*/
func (gog *Drive) BuildFileHierarchy(dirId string, buildPreview bool, includeFilter func(fileType string) bool) *Directory {

	//Get this item
	folderInfo, err := gog.connection.Files.
		Get(dirId).
		SupportsTeamDrives(true).
		Do()

	//Return nothing from this folder
	if err != nil {
		return nil
	}

	//Split the name and see if it is to be used
	name, date := gog.splitNameAndDate(folderInfo.Name)

	//Get all of the files in this folder
	dir := &Directory{
		File: File{
			Id:   folderInfo.Id,
			Name: name,
		},
		Type:  folderInfo.MimeType,
		Items: make([]Item, 0),
	}

	//If there is a date add it
	if date != nil {
		dir.Date = date
	}

	//Now get all of the files
	files, err := gog.connection.Files.List().
		SupportsTeamDrives(true).
		IncludeTeamDriveItems(true).
		Q("'" + dirId + "' in parents and trashed=false").
		Do()

	//If there is an error just return
	if err != nil {
		log.Printf("Unable to retrieve Drive client: %v\n", err)
		return nil
	}

	//For each file
	for _, item := range files.Files {
		//Make sure item is not trashed
		if !item.Trashed {
			//If the item is a folder, get all of it's children
			if item.MimeType == "application/vnd.google-apps.folder" {
				//Get the child
				childFolder := gog.BuildFileHierarchy(item.Id, buildPreview, includeFilter)

				//Now set the parent id to this
				childFolder.ParentId = dir.Id

				//Just add the child
				dir.Items = append(dir.Items, childFolder)

			} else if includeFilter(item.MimeType) { ////Else check the filter
				//Split the name and see if it is to be used
				name, date := gog.splitNameAndDate(item.Name)

				//Create a new document
				doc := &Document{
					File: File{
						Id:   item.Id,
						Name: name,
					},
					Type: item.MimeType,

					ParentId: dir.Id,
				}
				//Only build the previews if needed
				if buildPreview {
					doc.Preview = gog.GetFilePreview(item.Id)
					doc.ThumbnailUrl = gog.GetFileThumbnailUrl(item.Id)

				}

				//If there is a date add it
				if date != nil {
					doc.Date = date
				}

				//Store it
				dir.Items = append(dir.Items, doc)

			}
		}
	}

	return dir
}

/**
Builds all of the forms and downloads them at the same time
*/
func (gog *Drive) BuildFormHierarchy(dirId string) *Directory {

	//Get this item
	folderInfo, err := gog.connection.Files.
		Get(dirId).
		SupportsTeamDrives(true).
		Do()

	//Return nothing from this folder
	if err != nil {
		return nil
	}

	//Get all of the files in this folder
	dir := &Directory{
		File: File{
			Id:   folderInfo.Id,
			Name: folderInfo.Name,
		},
		Type:  folderInfo.MimeType,
		Items: make([]Item, 0),
	}

	//Now get all of the files
	files, err := gog.connection.Files.List().
		SupportsTeamDrives(true).
		IncludeTeamDriveItems(true).
		Q("'" + dirId + "' in parents").
		Do()

	//If there is an error just return
	if err != nil {
		log.Printf("Unable to retrieve Drive client: %v\n", err)
		return nil
	}

	//For each file
	for _, item := range files.Files {
		//Make sure item is not trashed
		if !item.Trashed {
			//If the item is a folder, get all of it's children
			if item.MimeType == "application/vnd.google-apps.folder" {
				//Get the child
				childFolder := gog.BuildFormHierarchy(item.Id)

				//Now set the parent id to this
				childFolder.ParentId = dir.Id

				//Only add the form is there are any children
				if len(childFolder.Items) > 0 {

					//Just add the child
					dir.Items = append(dir.Items, childFolder)
				}

			} else if item.MimeType == "application/json" {
				//Now download the forms
				form, err := gog.downloadForm(item.Id)

				//If there was an error
				if err != nil {
					log.Printf("Error: %v", err)
				} else {
					//Remove the extention
					name := strings.TrimSuffix(item.Name, filepath.Ext(item.Name))

					//Add the forms id
					form.Id = item.Id
					form.Name = name
					form.ParentId = dir.Id

					//Now add it to the parents children
					dir.Items = append(dir.Items, form)

				}

			}
		}

	}

	return dir
}

/**
* Method to get the information hierarchy
 */
func (gog *Drive) GetFilePreview(id string) string {
	//Get the file type
	fileInfo, err := gog.connection.Files.Get(id).SupportsTeamDrives(true).Do()
	if err != nil {
		log.Printf("Error: %v", err)
		return ""
	}

	//Only get the preview if it is a google doc
	if fileInfo.MimeType != "application/vnd.google-apps.document" {
		return ""
	}

	//Get the plain text version of the file
	resp, err := gog.connection.Files.Export(id, "text/plain").Download()

	//If there was an error just don't do anything
	if err != nil {
		log.Printf("Error: %v", err)
		return ""
	}
	//Get the entire thing
	result, _ := ioutil.ReadAll(resp.Body)

	//Remove extra white space
	space := regexp.MustCompile(`\s+`)
	resultString := space.ReplaceAllString(string(result), " ")

	//Return only the first specified number of chars
	//Get the minimum value
	previewLength := gog.previewLength
	if len(resultString) < previewLength {
		previewLength = len(resultString)
	}

	return resultString[0:previewLength]

}

/**
* Method to get the information hierarchy
 */
func (gog *Drive) downloadForm(id string) (*Form, error) {

	//Get the plain text version of the file
	resp, err := gog.connection.Files.Get(id).Download()

	//If there was an error just don't do anything
	if err != nil {
		return nil, err
	}

	//Encode the response
	dec := json.NewDecoder(resp.Body)

	//Createa a new forms
	form := &Form{}

	//Now decode the stream into the forms
	err = dec.Decode(form)

	//If there was an error just don't do anything
	if err != nil {
		return nil, err
	}

	//Return it
	return form, nil

}

/**
* Method to get the information hierarchy
 */
func (gog *Drive) GetFileThumbnailUrl(id string) string {

	//Get the file type
	fileInfo, err := gog.connection.Files.Get(id).SupportsTeamDrives(true).Do()
	if err != nil {
		log.Printf("Error: %v", err)
		return ""
	}

	//Only get the preview if it is a google doc
	if fileInfo.MimeType != "application/vnd.google-apps.document" {
		return ""
	}

	//Start up by getting the
	resp, err := gog.connection.Files.Export(id, "text/html").Download()

	//If there was an error just don't do anything
	if err != nil {
		log.Printf("Error: %v", err)
		return ""
	}
	//Now create a tokenizer
	tokenizer := html.NewTokenizer(resp.Body)

	for {
		//go to the next token
		tt := tokenizer.Next()

		if tt == html.ErrorToken {
			return ""
		} else if tt == html.StartTagToken {
			//Get the token
			t := tokenizer.Token()

			//If this is an image take
			if t.Data == "img" {
				//Now search for the source tag
				for _, a := range t.Attr {
					if a.Key == "src" {
						return a.Val
					}
				}
			}

		}

	}

}

/**
* Method to get the file html
 */
func (gog *Drive) GetFileHtml(id string) string {

	//Get the file type
	fileInfo, err := gog.connection.Files.Get(id).SupportsTeamDrives(true).Do()
	if err != nil {
		log.Printf("Error: %v", err)
		return ""
	}

	//If it is a pdf, get it is a pdf,
	switch fileInfo.MimeType {
	case "application/pdf":
		{

			//Get the plain text version of the file
			resp, err := gog.connection.Files.Get(id).Download()

			//If there was an error return
			if err != nil {
				return ""
			}
			defer resp.Body.Close()
			//Get the entire thing
			result, _ := ioutil.ReadAll(resp.Body)

			//Convert to base 64
			pdfBase64Str := base64.StdEncoding.EncodeToString(result)
			//Build the srcData
			srcData := "data:application/pdf;base64," + pdfBase64Str

			//Wrap in html
			html := "<embed  style=\"width:100%; height:80vh;\" type=\"application/pdf\" src=\"" + srcData + "\" />"
			html += "<a href=\"" + srcData + "\"> Open " + fileInfo.Name + " in full page view </a>"
			//Return only the first specified number of chars
			return html

		}
	default: //"application/vnd.google-apps.document"
		//Get the plain text version of the file
		resp, err := gog.connection.Files.Export(id, "text/html").Download()

		//If there was an error just don't do anything
		if err != nil {
			log.Printf("Error: %v", err)
			return ""
		}

		//Get the entire thing
		result, _ := ioutil.ReadAll(resp.Body)

		//Return only the first specified number of chars
		return string(result)
	}

}

/**
* Method to get the file html
 */
func (gog *Drive) GetArbitraryFile(id string) (io.ReadCloser, error) {

	//Get the plain text version of the file
	rep, err := gog.connection.Files.Get(id).Download()

	//If there was an error return
	if err != nil {
		return nil, err
	}

	//Ok return the read and closer
	return rep.Body, nil
}

/**
* Method to get the file html
 */
func (gog *Drive) GetMostRecentFileInDir(dirId string) (io.ReadCloser, error) {

	//Now get all of the files
	files, err := gog.connection.Files.List().
		SupportsTeamDrives(true).
		IncludeTeamDriveItems(true).
		Q("'" + dirId + "' in parents").
		OrderBy("recency desc").
		PageSize(1).
		Do()

	//If there is an error just return
	if err != nil {
		return nil, err
	}

	//There needs to be at least one file found
	if len(files.Files) < 1 {
		return nil, errors.New("no files not found in dir " + dirId)
	}

	//Get the plain text version of the file
	rep, err := gog.connection.Files.Get(files.Files[0].Id).Download()

	//If there was an error return
	if err != nil {
		return nil, err
	}

	//Ok return the read and closer
	return rep.Body, nil
}

/**
* Method to get the file html
 */
func (gog *Drive) GetFileAsInterface(id string, inter interface{}) error {
	//Get the resposne,
	rep, err := gog.GetArbitraryFile(id)
	defer rep.Close()
	//If there was no error
	if err != nil {
		return err
	}

	//REad the data
	data, err := ioutil.ReadAll(rep)
	if err != nil {
		return err
	}

	//Now decode the resposne into json
	err = json.Unmarshal(data, &inter)

	return err

}

/**
* Method to upload a file
 */
func (gog *Drive) PostArbitraryFile(fileName string, parent string, file io.Reader) (string, error) {
	//Create the file
	myFile := drive.File{
		Parents: []string{parent},
		Name:    fileName,
	}

	//Upload the file
	createdFile, err := gog.connection.Files.Create(&myFile).Media(file).SupportsTeamDrives(true).Do()
	if err != nil {
		return "", err
	}
	//Now return the link
	return createdFile.Id, nil

}

/**
  Gets the files matching the search in the dir

*/
func (gog *Drive) GetFirstFileMatching(dirId string, name string) (io.ReadCloser, error) {

	//Now get all of the files
	files, err := gog.connection.Files.List().
		SupportsTeamDrives(true).
		IncludeTeamDriveItems(true).
		Q("'" + dirId + "' in parents and trashed=false and fullText contains '" + name + "'").
		Do()

	if err != nil {
		return nil, err
	}

	//If there are no files
	if len(files.Files) < 1 {
		return nil, errors.New("no matching file found in for " + name)
	}

	//Now just return the file
	return gog.GetArbitraryFile(files.Files[0].Id)

}
