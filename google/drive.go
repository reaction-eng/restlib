// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package google

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/reaction-eng/restlib/configuration"
	"github.com/reaction-eng/restlib/file"
	"golang.org/x/net/context"
	"golang.org/x/net/html"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
)

type Drive struct {

	//Store the connection to google.  This has been wrapped with the correct Headers
	connection *drive.Service

	//Store the timezone
	timeZone string
}

//Get a new interface
func NewDrive(configuration configuration.Configuration) (*Drive, error) {
	//Create a new
	timeZone, err := configuration.GetStringError("default_time_zone")
	if err != nil {
		return nil, err
	}

	gInter := &Drive{
		timeZone: timeZone,
	}

	email, err := configuration.GetStringError("google_auth_email")
	if err != nil {
		return nil, err
	}
	privateKey, err := configuration.GetStringError("google_auth_key")
	if err != nil {
		return nil, err
	}

	//Open the client
	jwtConfig := &jwt.Config{
		Email:      email,
		PrivateKey: []byte(privateKey),
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

	return gInter, err
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

func (gog *Drive) BuildListing(dirId string, previewLength int, includeFilter func(fileType string) bool) (*file.Listing, error) {

	//Get this item
	folderInfo, err := gog.connection.Files.
		Get(dirId).
		SupportsTeamDrives(true).
		Do()

	//Return nothing from this folder
	if err != nil {
		return nil, err
	}

	//Split the name and see if it is to be used
	name, date := gog.splitNameAndDate(folderInfo.Name)

	//Get all of the files in this folder
	dir := file.NewListing()
	dir.Id = folderInfo.Id
	dir.Name = name

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
		return nil, err
	}

	//For each file
	for _, item := range files.Files {
		//Make sure item is not trashed
		if !item.Trashed {
			//If the item is a folder, get all of it's children
			if item.MimeType == "application/vnd.google-apps.folder" {
				//Get the child
				childFolder, err := gog.BuildListing(item.Id, previewLength, includeFilter)

				if err != nil {
					return nil, err
				}

				//Now set the parent Id to this
				childFolder.ParentId = dir.Id

				//Just add the child
				dir.Listings = append(dir.Listings, *childFolder)

			} else if includeFilter(item.MimeType) { ////Else check the filter
				//Split the name and see if it is to be used
				name, date := gog.splitNameAndDate(item.Name)

				//Create a new document
				doc := file.Item{
					Id:       item.Id,
					Name:     name,
					ParentId: dir.Id,
					Type:     item.MimeType,
				}

				//Only build the previews if needed
				if previewLength > 0 {
					doc.Preview = gog.GetFilePreview(item.Id, previewLength)
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

	return dir, nil
}

/**
* Method to get the information hierarchy
 */
func (gog *Drive) GetFilePreview(id string, previewLength int) string {
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
	if len(resultString) < previewLength {
		previewLength = len(resultString)
	}

	return resultString[0:previewLength]

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
	//If there was no error
	if err != nil {
		return err
	}
	defer rep.Close()
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
func (gog *Drive) PostArbitraryFile(fileName string, parent string, file io.Reader, mime string) (string, error) {
	//Create the file
	myFile := drive.File{
		Parents: []string{parent},
		Name:    fileName,
	}

	//If there is a mime type use it
	if len(mime) > 0 {
		myFile.MimeType = mime
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
