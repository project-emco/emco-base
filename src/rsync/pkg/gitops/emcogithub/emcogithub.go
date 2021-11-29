package emcogithub

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fluxcd/go-git-providers/github"
	"github.com/fluxcd/go-git-providers/gitprovider"
)

const (
	githubDomain = "github.com"
)

/*
	Function to create gitprovider githubClient
	params : github token
	return : gitprovider github client, error
*/
func CreateClient(githubToken string) (gitprovider.Client, error) {
	c, err := github.NewClient(gitprovider.WithOAuth2Token(githubToken), gitprovider.WithDestructiveAPICalls(true))
	if err != nil {
		return nil, err
	}
	return c, nil
}

/*
	Function to create a new Repo in github
	params : context, github client, Repository Name, User Name, description
	return : nil/error
*/

func CreateRepo(ctx context.Context, c gitprovider.Client, repoName string, userName string, desc string) error {

	// create repo reference
	userRepoRef := getRepoRef(userName, repoName)

	// Create repoinfo reference
	userRepoInfo := gitprovider.RepositoryInfo{
		Description: &desc,
		Visibility:  gitprovider.RepositoryVisibilityVar(gitprovider.RepositoryVisibilityPublic),
	}

	// Create the repository
	_, err := c.UserRepositories().Create(ctx, userRepoRef, userRepoInfo, &gitprovider.RepositoryCreateOptions{
		AutoInit:        gitprovider.BoolVar(true),
		LicenseTemplate: gitprovider.LicenseTemplateVar(gitprovider.LicenseTemplateApache2),
	})

	if err != nil {
		return err
	}
	fmt.Println("Repo Created")

	return nil
}

/*
iterate through all the files and create a file array with path and contents
params : path string
return : []gitprovider.CommitFile
*/
func Iterate(path string) []gitprovider.CommitFile {
	files := []gitprovider.CommitFile{}
	length := len([]rune(path))
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		// check if the path is a directory
		if !info.IsDir() {

			fmt.Printf("File Name: %s\n", info.Name())
			//Get the contents of the file
			content := GetContent(path)
			//get the relative path of the file
			relativePath := path[length+1:]
			fmt.Println(relativePath)
			// append to array of files
			files = append(files, gitprovider.CommitFile{Path: &relativePath, Content: &content})
		}
		return nil
	})
	fmt.Println(len(files))
	return files
}

/*
	Function to get the contents of the file
	params : filepath (string)
	return : contents (string)
*/
func GetContent(filePath string) string {

	// Open file on disk.
	f, _ := os.Open(filePath)

	// Read entire file into byte slice.
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	return string(content)
}

/*
	Function to commit multiple files to the github repo
	params : context, github client, User Name, Repo Name, Branch Name, Commit Message, files ([]gitprovider.CommitFile)
	return : nil/error
*/

func CommitFiles(ctx context.Context, c gitprovider.Client, userName string, repoName string, branch string, commitMessage string, files []gitprovider.CommitFile) error {

	// create repo reference
	userRepoRef := getRepoRef(userName, repoName)

	userRepo, err := c.UserRepositories().Get(ctx, userRepoRef)
	if err != nil {
		return err
	}
	//Commit file to this repo
	_, err = userRepo.Commits().Create(ctx, branch, commitMessage, files)

	if err != nil {
		return err
	}
	return nil
}

/*
	Function to delete repo
	params : context, gitprovider client , user name, repo name
	return : nil/error
*/
func DeleteRepo(ctx context.Context, c gitprovider.Client, userName string, repoName string) error {

	// create repo reference
	userRepoRef := getRepoRef(userName, repoName)
	// get the reference of the repo to be deleted
	userRepo, err := c.UserRepositories().Get(ctx, userRepoRef)

	if err != nil {
		return err
	}
	//delete repo
	fmt.Println(userRepo.Delete(ctx))

	return nil
}

/*
	Internal function to create a repo refercnce
	params : user name, repo name
	return : repo reference
*/

func getRepoRef(userName string, repoName string) gitprovider.UserRepositoryRef {
	// Create the user reference
	userRef := gitprovider.UserRef{
		Domain:    githubDomain,
		UserLogin: userName,
	}

	// Create the repo reference
	userRepoRef := gitprovider.UserRepositoryRef{
		UserRef:        userRef,
		RepositoryName: repoName,
	}

	return userRepoRef
}

/*
	Function to Add file to the commit
	params : path , content, files (gitprovider commitfile array)
	return : files (gitprovider commitfile array)
*/
func Add(path string, content string, files []gitprovider.CommitFile) []gitprovider.CommitFile {
	files = append(files, gitprovider.CommitFile{
		Path:    &path,
		Content: &content,
	})

	return files
}

/*
	Function to Delete file from the commit
	params : path, files (gitprovider commitfile array)
	return : files (gitprovider commitfile array)
*/
func Delete(path string, files []gitprovider.CommitFile) []gitprovider.CommitFile {
	files = append(files, gitprovider.CommitFile{
		Path:    &path,
		Content: nil,
	})

	return files
}