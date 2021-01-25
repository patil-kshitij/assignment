package repositories

import(
	"errors"
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"

	"assignment/httprest"
	"assignment/config"
	"assignment/entity"
	"encoding/base64"

)

const(
	createBranchURL = "%s/repos/%s/%s/git/refs"
	ContentURL = "%s/repos/%s/%s/contents/%s"
	PullReqURL = "%s/repos/%s/%s/pulls"
	SHAReqURL = "%s/repos/%s/%s/branches/%s"
	//createOrUpdateFileContentURL = "%s/repos/%s/%s/contents/%s"
	
	masterBranch = "master"

	refConst = "refs/heads/%s"

	authHeaderKey = "Authorization"
	AuthToken = "token %s"

	AcceptKey = "Accept"
	GithubAcceptValue = "application/vnd.github.v3+json"


)

func CreateGithubImpl()GithubImpl{
	return GithubImpl{}
}

type GithubImpl struct{}

func (g GithubImpl) GetMasterSHA(owner,repo,token string) (string,error){
	url := fmt.Sprintf(SHAReqURL,config.AppConfig.GithubAPIURL,owner,repo,masterBranch)

	headers := make(map[string]string)
	headers[authHeaderKey] = fmt.Sprintf(AuthToken, token)
	headers[AcceptKey] = GithubAcceptValue

	resp,err := httprest.CallRestAPI(http.MethodGet,url,headers,nil)
	if err != nil {
		return "",err
	}
	if resp.StatusCode != http.StatusOK {
		return "",errors.New("File Content Not obtained")
	}

	apiOutput,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "",err
	}

	shaResp := &entity.SHAResp{}
	err = json.Unmarshal(apiOutput,shaResp)
	if err!= nil {
		return "",err
	}

	return shaResp.Commit.Sha,nil


}

func (g GithubImpl) CreateBranch(owner, repo, branchName, refShaID,token string) error {
	url := fmt.Sprintf(createBranchURL,config.AppConfig.GithubAPIURL,owner,repo)
	
	headers := make(map[string]string)
	headers[authHeaderKey] = fmt.Sprintf(AuthToken, token)
	
	reqBody := entity.BranchReq{
		Ref: fmt.Sprintf(refConst,branchName),
		Sha: refShaID,
	}
	reqJSON,err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	resp,err := httprest.CallRestAPI(http.MethodPost,url,headers,reqJSON)
	if err != nil {
		return err
	}

	apiOutput,err := ioutil.ReadAll(resp.Body)
	fmt.Println("ApiOutput branch = ",string(apiOutput))
	fmt.Println("Error in reading branch create output",err)
	

	if resp.StatusCode != http.StatusCreated {
		fmt.Println("Status = ",resp.StatusCode)
		return errors.New("Branch not created")
	}
	return nil
}

func (g GithubImpl) GetFileContent(owner,repo,branchName,fileName,token string) (*entity.FileContentResp,error) {
	url := fmt.Sprintf(ContentURL,config.AppConfig.GithubAPIURL,owner,repo,fileName)

	headers := make(map[string]string)
	headers[authHeaderKey] = fmt.Sprintf(AuthToken, token)
	resp,err := httprest.CallRestAPI(http.MethodGet,url,headers,nil)
	if err != nil {
		return nil,err
	}
	if resp.StatusCode != http.StatusOK {
		return nil,errors.New("File Content Not obtained")
	}
	apiOutput,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,err
	}
	fileContent := &entity.FileContentResp{}
	err = json.Unmarshal(apiOutput,fileContent)
	if err != nil {
		return nil,err
	}

	Content,err := base64.StdEncoding.DecodeString(fileContent.Content)
	if err != nil {
		return nil,err
	}
	fileContent.Content = string(Content)

	return fileContent,nil
}

func (g GithubImpl) CreateORUpdateFile(owner,repo,fileName,token string, req *entity.UpdateFileReq)error{
	url := fmt.Sprintf(ContentURL,config.AppConfig.GithubAPIURL,owner,repo,fileName)

	headers := make(map[string]string)
	headers[authHeaderKey] = fmt.Sprintf(AuthToken, token)
	headers[AcceptKey] = GithubAcceptValue

	reqJSON,err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp,err := httprest.CallRestAPI(http.MethodPut,url,headers,reqJSON)
	if err != nil {
		return err
	}

	if  resp.StatusCode != http.StatusOK || resp.StatusCode != http.StatusCreated {
		return err
	}

	return nil
}

func (g GithubImpl) CreatePullRequest(owner, repo, token string, req entity.PullReq) error {
	url := fmt.Sprintf(PullReqURL,config.AppConfig.GithubAPIURL,owner,repo)
	
	headers := make(map[string]string)
	headers[authHeaderKey] = fmt.Sprintf(AuthToken, token)
	headers[AcceptKey] = GithubAcceptValue

	reqJSON,err := json.Marshal(req)
	fmt.Println("Create Pull request json marshal err =",err)
	if err != nil {
		return err
	}

	fmt.Println("Pull Req JSON",string(reqJSON))

	resp,err := httprest.CallRestAPI(http.MethodPost,url,headers,reqJSON)
	fmt.Println("Create Pull request call err =",err)
	if err != nil {
		return err
	}

	apiOutput,err := ioutil.ReadAll(resp.Body)
	fmt.Println("ApiOutput create PR = ",string(apiOutput))
	fmt.Println("Error in reading PR create output",err)

	if  resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errors.New("Some erroneous status")
	}

	return nil
}
