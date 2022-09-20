// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package controllers

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	k8spluginv1alpha1 "gitlab.com/project-emco/core/emco-base/src/monitor/pkg/apis/k8splugin/v1alpha1"
	gitsupport "gitlab.com/project-emco/core/emco-base/src/monitor/pkg/gitops/gitsupport"
	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
)

const (
	maxrand = 0x7fffffffffffffff
)

var mutex = sync.Mutex{}

// func getCredCallBack(userName, token string) func(url string, username string, allowedTypes git.CredType) (*git.Credential, error) {
// 	return func(url string, username string, allowedTypes git.CredType) (*git.Credential, error) {
// 		username = userName
// 		password := token
// 		cred, err := git.NewCredUserpassPlaintext(username, password)
// 		return cred, err
// 	}
// }

// CommitFile contains high-level information about a file added to a commit.
// type CommitFile struct {
// 	// Path is path where this file is located.
// 	// +required
// 	Add bool `json:"add"`

// 	Path *string `json:"path"`

// 	FileName *string `json:"filename"`

// 	// Content is the content of the file.
// 	// +required
// 	Content *string `json:"content,omitempty"`
// }

// type GitClient struct {
// 	gitProviderClient gitprovider.Client
// 	gogithubClient    *gogithub.Client
// }

type GithubAccessClient struct {
	gitUser string
	gitRepo string
	cluster string

	gitProvider gitsupport.GitProvider
}

var GitHubClient GithubAccessClient

// /*
// 	Function to create gitClient
// 	params : userName, github token
// 	return : github client, error
// */
// func CreateClient(userName, githubToken string) (GitClient, error) {

// 	var client GitClient
// 	var err error

// 	client.gitProviderClient, err = github.NewClient(gitprovider.WithOAuth2Token(githubToken), gitprovider.WithDestructiveAPICalls(true))
// 	if err != nil {
// 		return GitClient{}, err
// 	}

// 	tp := gogithub.BasicAuthTransport{
// 		Username: userName,
// 		Password: githubToken,
// 	}
// 	client.gogithubClient = gogithub.NewClient(tp.Client())

// 	return client, nil

// }

func SetupGitHubClient() error {
	var err error
	GitHubClient, err = NewGitHubClient()
	return err
}

func NewGitHubClient() (GithubAccessClient, error) {

	gitUser := os.Getenv("GIT_USERNAME")
	gitToken := os.Getenv("GIT_TOKEN")
	gitRepo := os.Getenv("GIT_REPO")
	clusterName := os.Getenv("GIT_CLUSTERNAME")

	// If any value is not provided then can't store in Git location
	if len(gitRepo) <= 0 || len(gitToken) <= 0 || len(gitUser) <= 0 || len(clusterName) <= 0 {
		log.Info("Github information not found:: Skipping Github storage", log.Fields{})
		return GithubAccessClient{}, nil
	}
	log.Info("GitHub Info found", log.Fields{"gitRepo::": gitRepo, "cluster::": clusterName})

	// cl, err := CreateClient(gitUser, gitToken)
	// if err != nil {
	// 	return GithubAccessClient{}, err
	// }

	gitProvider, err := gitsupport.NewGitProvider()
	if err != nil {
		return GithubAccessClient{}, err
	}

	p := GithubAccessClient{
		gitUser:     gitUser,
		gitRepo:     gitRepo,
		cluster:     clusterName,
		gitProvider: *gitProvider,
	}

	return p, nil
}

func (c *GithubAccessClient) CommitCRToGitHub(cr *k8spluginv1alpha1.ResourceBundleState, l map[string]string) error {

	//Check if Github Client is available
	// if c.cl == (GitClient{}) {
	// 	return nil
	// }
	resBytes, err := json.Marshal(cr)
	if err != nil {
		log.Info("json Marshal error for resource::", log.Fields{"cr": cr, "err": err})
		return err
	}
	// Get cid and app id
	v, ok := l["emco/deployment-id"]
	if !ok {
		return fmt.Errorf("Unexpected error:: Inconsistent labels %v", l)
	}
	result := strings.SplitN(v, "-", 2)
	if len(result) != 2 {
		return fmt.Errorf("Unexpected error:: Inconsistent labels %v", l)
	}
	app := result[1]
	cid := result[0]
	path := "clusters/" + c.cluster + "/status/" + cid + "/app/" + app + "/" + v

	// Add files for commit

	folderName := "/tmp/" + c.gitUser + "-" + c.gitRepo + "-" + c.cluster
	var files interface{}
	files, err = c.gitProvider.Apply(path, files, resBytes)
	if err != nil {
		log.Error("Error in Applying files", log.Fields{"err": err, "files": files, "path": path})
		return err
	}
	branchName := c.cluster

	//commit files
	commitMessage := "Adding Status for " + path + " to branch " + branchName

	// commitfiles
	// err = commitFiles(c.url, c.gitToken, c.gitUser, commitMessage, appName, folderName, "main", files)
	mutex.Lock()
	defer mutex.Unlock()
	err = c.gitProvider.CommitFiles(commitMessage, branchName, folderName, files)
	if err != nil {
		log.Error("ApplyConfig:: Commit files err", log.Fields{"err": err, "files": files})
		return err
	}

	return nil
}

// /*
// 	Internal function to create a repo refercnce
// 	params : user name, repo name
// 	return : repo reference
// */
// func (c *GithubAccessClient) getRepoRef(userName string, repoName string) gitprovider.UserRepositoryRef {
// 	// Create the user reference
// 	userRef := gitprovider.UserRef{
// 		Domain:    c.gitType + ".com",
// 		UserLogin: userName,
// 	}

// 	// Create the repo reference
// 	userRepoRef := gitprovider.UserRepositoryRef{
// 		UserRef:        userRef,
// 		RepositoryName: repoName,
// 	}

// 	return userRepoRef
// }

// var mutex = sync.Mutex{}

// /*
// 	Function to commit multiple files to the github repo
// 	params : context, Branch Name, Commit Message, appName, files ([]gitprovider.CommitFile)
// 	return : nil/error
// */
// func (c *GithubAccessClient) CommitFiles(ctx context.Context, branch, commitMessage, appName string, files []gitprovider.CommitFile) error {

// 	mergeBranch := appName
// 	// Only one process to commit to Github location to avoid conflicts
// 	mutex.Lock()
// 	defer mutex.Unlock()

// 	// commit the files to this new branch
// 	// create repo reference
// 	log.Info("Creating Repo Reference. ", log.Fields{})
// 	userRepoRef := c.getRepoRef(c.gitUser, c.gitRepo)
// 	log.Info("UserRepoRef:", log.Fields{"UserRepoRef": userRepoRef})

// 	log.Info("Obtaining user repo. ", log.Fields{})
// 	userRepo, err := c.cl.gitProviderClient.UserRepositories().Get(ctx, userRepoRef)
// 	if err != nil {
// 		log.Error("Error in commiting the files", log.Fields{"err": err, "mergeBranch": mergeBranch, "commitMessage": commitMessage, "files": files})
// 		return err
// 	}
// 	log.Info("UserRepo:", log.Fields{"UserRepo": userRepo})

// 	log.Info("Commiting Files:", log.Fields{"files": files})
// 	//Commit file to this repo
// 	resp, err := userRepo.Commits().Create(ctx, mergeBranch, commitMessage, files)
// 	if err != nil {
// 		if !strings.Contains(err.Error(), "404 Not Found") {
// 			log.Error("Error in commiting the files", log.Fields{"err": err, "mergeBranch": mergeBranch, "commitMessage": commitMessage, "files": files})
// 		}
// 		return err

// 	}
// 	log.Info("CommitResponse for userRepo:", log.Fields{"resp": resp})
// 	return nil
// }

// /*
// 	Function to obtaion the SHA of latest commit
// 	params : context, git client, User Name, Repo Name, Branch, Path
// 	return : LatestCommit string, error
// */
// func GetLatestCommitSHA(ctx context.Context, c GitClient, userName, repoName, branch, path string) (string, error) {

// 	perPage := 1
// 	page := 1

// 	lcOpts := &gogithub.CommitsListOptions{
// 		ListOptions: gogithub.ListOptions{
// 			PerPage: perPage,
// 			Page:    page,
// 		},
// 		SHA:  branch,
// 		Path: path,
// 	}
// 	//Get the latest SHA
// 	resp, _, err := c.gogithubClient.Repositories.ListCommits(ctx, userName, repoName, lcOpts)
// 	if err != nil {
// 		log.Error("Error in obtaining the list of commits", log.Fields{"err": err})
// 		return "", err
// 	}
// 	if len(resp) == 0 {
// 		log.Debug("File not created yet.", log.Fields{"Latest Commit Array": resp})
// 		return "", nil
// 	}
// 	latestCommitSHA := *resp[0].SHA

// 	return latestCommitSHA, nil
// }

// /*
// 	Function to create a new branch
// 	params : context, git client,latestCommitSHA, User Name, Repo Name, Branch
// 	return : error
// */
// func createBranch(ctx context.Context, c GitClient, latestCommitSHA, userName, repoName, branch string) error {
// 	// create a new branch
// 	ref, _, err := c.gogithubClient.Git.CreateRef(ctx, userName, repoName, &gogithub.Reference{
// 		Ref: gogithub.String("refs/heads/" + branch),
// 		Object: &gogithub.GitObject{
// 			SHA: gogithub.String(latestCommitSHA),
// 		},
// 	})
// 	if err != nil {
// 		log.Error("Git.CreateRef returned error:", log.Fields{"err": err})
// 		return err

// 	}
// 	log.Info("Branch Created: ", log.Fields{"ref": ref})
// 	return nil
// }

// //function to delete a file
// func deleteFile(filenName string) error {
// 	// Removing file from the directory
// 	// Using Remove() function
// 	err := os.Remove(filenName)
// 	if err != nil {
// 		log.Error("Error in Deleting file from the tmp folder", log.Fields{"err": err})
// 		return err
// 	}

// 	return nil
// }

// //function to create a new file
// func createFile(fileName string, content string) error {
// 	if err := os.MkdirAll(filepath.Dir(fileName), 0770); err != nil {
// 		return err
// 	}

// 	f, err := os.Create(fileName)

// 	if err != nil {
// 		log.Error("Error in Creating file in the tmp folder", log.Fields{"err": err})
// 		return err
// 	}

// 	defer f.Close()

// 	_, err2 := f.WriteString(content)

// 	if err2 != nil {
// 		log.Error("Error in writing file from the tmp folder", log.Fields{"err": err2})
// 		return err2
// 	}

// 	return nil
// }

// // function to commit files to a branch
// func commitFiles(url, token, userName, commitMessage, appName, folderName, branch string, files []CommitFile) error {

// 	// Only one process to commit to Github location to avoid conflicts
// 	mutex.Lock()
// 	defer mutex.Unlock()

// 	branchName := appName

// 	// clone git the repo to local repo
// 	check, err := exists(folderName)

// 	if !check {
// 		if err := os.Mkdir(folderName, os.ModePerm); err != nil {
// 			log.Error("Error in creating the dir", log.Fields{"Error": err})
// 			return err
// 		}
// 		// // clone the repo
// 		log.Info("URL", log.Fields{"URL": url})
// 		fmt.Println("URL %s", url)
// 		repo, err := git.Clone(url, folderName, &git.CloneOptions{CheckoutBranch: branch, CheckoutOptions: git.CheckoutOptions{Strategy: git.CheckoutSafe}})
// 		if err != nil {
// 			log.Error("Error cloning the repo", log.Fields{"Error": err})
// 			return err
// 		}
// 		fmt.Println(repo)
// 	}

// 	// // // open a repo
// 	repo, err := git.OpenRepository(folderName)
// 	if err != nil {
// 		log.Error("Error in Opening the git repository", log.Fields{"err": err, "appName": appName})
// 		return err
// 	}

// 	signature := &git.Signature{
// 		Name:  userName,
// 		Email: userName + "@gmail.com",
// 		When:  time.Now(),
// 	}

// 	var targetID *git.Oid
// 	//create a new branch
// 	//check if branch already exists then skip create branch
// 	// check if a local branch exitst if not do a checkout
// 	localBranch, err := repo.LookupBranch(branchName, git.BranchLocal)
// 	// No local branch, lets create one
// 	if localBranch == nil || err != nil {
// 		branchHandle, err := CreateBranch(folderName, branchName)
// 		if err != nil {
// 			if !strings.Contains(err.Error(), "a reference with that name already exists") {
// 				return err
// 			}
// 		}
// 		targetID = branchHandle.Target()
// 	} else {
// 		branchRef, err := repo.References.Lookup("refs/heads/" + branchName)
// 		if err != nil {
// 			log.Info("Error in looking up ref", log.Fields{"err": err})
// 			return err
// 		}

// 		targetID = branchRef.Target()
// 	}

// 	//commit files to the branch
// 	//push the branch

// 	// set head to point to the created branch
// 	err = repo.SetHead("refs/heads/" + branchName)
// 	if err != nil {
// 		log.Error("Error in settting the head", log.Fields{"err": err, "branchName": branchName})
// 		return err
// 	}

// 	//Update the index with files and obtain the latest index
// 	//loop through all files and update the index
// 	idx, err := repo.Index()
// 	if err != nil {
// 		log.Error("Error in obtaining the repo index", log.Fields{"err": err, "idx": idx})
// 		return err
// 	}

// 	for _, file := range files {
// 		if file.Add {
// 			idx, err = addToCommit(idx, *file.Path, *file.FileName, *file.Content)
// 		} else {
// 			idx, err = deleteFromCommit(idx, *file.Path, *file.FileName)
// 		}

// 		if err != nil {
// 			log.Error("Error in adding or deleting file to commit", log.Fields{"err": err, "idx": idx})
// 			return err
// 		}
// 	}
// 	//commit the files to the branch
// 	treeId, err := idx.WriteTree()
// 	if err != nil {
// 		log.Error("Error from idx.WriteTree()", log.Fields{"err": err})
// 		return err
// 	}

// 	err = idx.Write()
// 	if err != nil {
// 		log.Error("Error in Deleting file from idx.Write()", log.Fields{"err": err})
// 		return err
// 	}

// 	tree, err := repo.LookupTree(treeId)
// 	if err != nil {
// 		log.Error("Error in looking up tree", log.Fields{"err": err, "treeId": treeId})
// 		return err
// 	}

// 	commitTarget, err := repo.LookupCommit(targetID)
// 	if err != nil {
// 		log.Error("Error in Looking up Commit for commit", log.Fields{"err": err})
// 		return err
// 	}

// 	_, err = repo.CreateCommit("refs/heads/"+branchName, signature, signature, commitMessage, tree, commitTarget)
// 	if err != nil {
// 		log.Error("Error in creating a commit", log.Fields{"err": err, "branchName": branchName})
// 		return err
// 	}
// 	err = pushBranch(repo, branchName, userName, token)

// 	return nil
// }

// // function to push branch to remote origin
// func pushBranch(repo *git.Repository, branchName, userName, token string) error {
// 	// push the branch to origin
// 	remote, err := repo.Remotes.Lookup("origin")

// 	cbs := &git.RemoteCallbacks{
// 		CredentialsCallback: getCredCallBack(userName, token),
// 	}

// 	err = remote.Push([]string{"+refs/heads/" + branchName}, &git.PushOptions{RemoteCallbacks: *cbs})
// 	if err != nil {
// 		log.Error("Error in Pushing the branch", log.Fields{"err": err, "branchName": branchName})
// 		return err
// 	}

// 	return nil
// }

// // function to push branch to remote origin
// func pushDeleteBranch(repo *git.Repository, branchName, userName, token string) error {
// 	// push the branch to origin
// 	remote, err := repo.Remotes.Lookup("origin")

// 	cbs := &git.RemoteCallbacks{
// 		CredentialsCallback: getCredCallBack(userName, token),
// 	}

// 	err = remote.Push([]string{":refs/heads/" + branchName}, &git.PushOptions{RemoteCallbacks: *cbs})
// 	if err != nil {
// 		log.Error("Error in Pushing the branch", log.Fields{"err": err, "branchName": branchName})
// 		return err
// 	}

// 	return nil
// }

// //function to merge branch to main (Should include a commit as well)
// func mergeToMain(repo *git.Repository, branchName string, signature *git.Signature) error {
// 	// get reference for the target merge branch
// 	mergeBranch, err := repo.References.Lookup("refs/heads/" + branchName)
// 	if err != nil {
// 		log.Error("Error in obtaining the reference for branch to merge to main", log.Fields{"err": err, "branchName": branchName})
// 		return err
// 	}

// 	mergeHeadMergeBranch, err := repo.AnnotatedCommitFromRef(mergeBranch)
// 	if err != nil {
// 		log.Error("Error in obtaining the head of the branch to merge", log.Fields{"err": err, "mergeHeadMergeBranch": mergeHeadMergeBranch})
// 		return err
// 	}
// 	mergeHeads := make([]*git.AnnotatedCommit, 1)

// 	mergeHeads[0] = mergeHeadMergeBranch

// 	err = repo.Merge(mergeHeads, nil, nil)
// 	if err != nil {
// 		log.Error("Error in Merging the branch", log.Fields{"err": err, "mergeHeads": mergeHeads})
// 		return err
// 	}

// 	mergeMessage, err := repo.Message()
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println("Merge Message: ", mergeMessage)

// 	err = commitMergeToMaster(repo, signature, "Merge commit to main")
// 	if err != nil {
// 		log.Error("Error in commit Merge to main", log.Fields{"err": err})
// 		return err
// 	}

// 	return nil
// }

// //function to delete branch
// func deleteBranch(repo *git.Repository, branchName string) error {
// 	branchA, err := repo.LookupBranch(branchName, git.BranchLocal)
// 	err = branchA.Delete()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// //function to check if folder exists
// // exists returns whether the given file or directory exists
// func exists(path string) (bool, error) {
// 	_, err := os.Stat(path)
// 	if err == nil {
// 		return true, nil
// 	}
// 	if os.IsNotExist(err) {
// 		return false, nil
// 	}
// 	return false, err
// }

// // function to add file for commit
// func addToCommit(idx *git.Index, path, fileName, contents string) (*git.Index, error) {

// 	err := createFile(path, contents)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// add file to staging area
// 	err = idx.AddByPath(fileName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return idx, nil

// }

// // function to delete file for commit
// func deleteFromCommit(idx *git.Index, path, fileName string) (*git.Index, error) {
// 	err := idx.RemoveByPath(fileName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = deleteFile(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return idx, nil
// }

// // function to add file to commit files array
// func add(path, fileName, content string, ref interface{}) []CommitFile {
// 	files := append(convertToCommitFile(ref), CommitFile{
// 		Add:      true,
// 		Path:     &path,
// 		FileName: &fileName,
// 		Content:  &content,
// 	})

// 	return files

// }

// // function to delete file from commit files array
// func delete(path, fileName string, ref interface{}) []CommitFile {
// 	files := append(convertToCommitFile(ref), CommitFile{
// 		Add:      false,
// 		Path:     &path,
// 		FileName: &fileName,
// 	})

// 	return files
// }

// func commitMergeToMaster(repo *git.Repository, signature *git.Signature, message string) error {
// 	//commit the merge to main
// 	branchName := "main"
// 	idx, err := repo.Index()
// 	if err != nil {
// 		log.Error("commitMergeToMaster: Error in obtaining the repo index", log.Fields{"err": err, "idx": idx})
// 		return err
// 	}

// 	branchMain, err := repo.LookupBranch(branchName, git.BranchLocal)
// 	if err != nil {
// 		return err
// 	}

// 	treeId, err := idx.WriteTree()
// 	if err != nil {
// 		log.Error("commitMergeToMaster: Error from idx.WriteTree()", log.Fields{"err": err})
// 		return err
// 	}

// 	err = idx.Write()
// 	if err != nil {
// 		log.Error("commitMergeToMaster: Error in Deleting file from idx.Write()", log.Fields{"err": err})
// 		return err
// 	}

// 	tree, err := repo.LookupTree(treeId)
// 	if err != nil {
// 		log.Error("commitMergeToMaster: Error in looking up tree", log.Fields{"err": err, "treeId": treeId})
// 		return err
// 	}

// 	commitTarget, err := repo.LookupCommit(branchMain.Target())
// 	if err != nil {
// 		log.Error("commitMergeToMaster: Error in Looking up Commit for commit", log.Fields{"err": err})
// 		return err
// 	}

// 	_, err = repo.CreateCommit("refs/heads/"+branchName, signature, signature, message, tree, commitTarget)
// 	if err != nil {
// 		log.Error("commitMergeToMaster:Error in creating a commit", log.Fields{"err": err, "branchName": branchName})
// 		return err
// 	}

// 	return nil
// }

// /*
// 	Helper function to convert interface to []gitprovider.CommitFile
// 	params: files interface{}
// 	return: []gitprovider.CommitFile
// */
// func convertToCommitFile(ref interface{}) []CommitFile {
// 	var exists bool
// 	switch ref.(type) {
// 	case []CommitFile:
// 		exists = true
// 	default:
// 		exists = false
// 	}
// 	var rf []CommitFile
// 	// Create rf is doesn't exist
// 	if !exists {
// 		rf = []CommitFile{}
// 	} else {
// 		rf = ref.([]CommitFile)
// 	}
// 	return rf
// }

// //function to create branch
// func CreateBranch(folderName, branchName string) (*git.Branch, error) {
// 	// // // open a repo
// 	repo, err := git.OpenRepository(folderName)
// 	if err != nil {
// 		log.Error("Error in Opening the git repository", log.Fields{"err": err})
// 		return nil, err
// 	}

// 	// create the new branch (May cause problems, try to get the headCommit of main)
// 	//checkout the new branch
// 	// set head to point to the created branch
// 	err = repo.SetHead("refs/heads/" + "main")
// 	if err != nil {
// 		log.Error("Error in settting the head", log.Fields{"err": err, "branchName": branchName})
// 		return nil, err
// 	}

// 	head, err := repo.Head()
// 	if err != nil {
// 		log.Error("Error in obtaining the head of the repo", log.Fields{"err": err})
// 		return nil, err
// 	}

// 	headCommit, err := repo.LookupCommit(head.Target())
// 	if err != nil {
// 		log.Error("Error in obtainging the head commit", log.Fields{"err": err, "headCommit": headCommit})
// 		return nil, err
// 	}
// 	branch, err := repo.CreateBranch(branchName, headCommit, false)
// 	if err != nil {
// 		log.Error("Error in Creating branch", log.Fields{"err": err, "branchName": branchName, "headCommit": headCommit, "branch": branch})
// 		return nil, err
// 	}

// 	return branch, nil
// }

//function to delete status folder for git
func (c *GithubAccessClient) DeleteStatusFromGit(appName string) error {

	// s := strings.SplitN(appName, "-", 2)
	// cid := s[0]
	// app := s[1]
	// path := "clusters/" + c.cluster + "/status/" + cid + "/app/" + app + "/" + appName
	// folderName := "/tmp/" + c.gitUser + "-" + c.gitRepo + "-" + c.cluster
	// statusBranchName := c.cluster

	// //get files to be deleted
	// // files, err := getFilesToDelete(folderName, path)
	// // if err != nil {
	// // 	log.Error("Error in obtaining files to be deleted", log.Fields{"folderName": folderName, "path": path})
	// // }

	// files := delete(folderName+"/"+path, path, []CommitFile{})

	// err := commitFiles(c.url, c.gitToken, c.gitUser, "Deleting status for "+appName, statusBranchName, folderName, "main", files)

	// return err

	return nil
}

// //function to delete all files in a path
// func getFilesToDelete(folderName, filePath string) ([]CommitFile, error) {
// 	var files []CommitFile
// 	path := folderName + "/" + filePath
// 	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
// 		file, err := os.Open(path)
// 		fileInfo, err := file.Stat()
// 		if err != nil {
// 			return err
// 		}

// 		//Add only if the path is not a directory
// 		if !fileInfo.IsDir() {
// 			s := strings.TrimPrefix(path, folderName+"/")
// 			files = delete(path, s, files)
// 		}
// 		return nil
// 	})

// 	return files, err
// }
