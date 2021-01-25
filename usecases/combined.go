package usecases

import(
	"context"
	"strings"
	"fmt"
	"assignment/constants"
	"encoding/base64"
	"errors"
	"assignment/entity"
)

type Github interface{
	GetMasterSHA(owner,repo,token string) (string,error)
	CreateBranch(owner, repo, branchName, refShaID,token string) error
	GetFileContent(owner,repo,branchName,fileName,token string) (*entity.FileContentResp,error)
	CreateORUpdateFile(owner,repo,fileName,token string, req *entity.UpdateFileReq)error
	CreatePullRequest(owner, repo, token string, req entity.PullReq) error
}


func CreateCombinedUseCase(repo Github) CombinedUseCase {
	return CombinedUseCase{
		githubRepo:repo,
	}
}

type CombinedUseCase struct {
	githubRepo Github
}

func (cuc CombinedUseCase) Perform(ctx context.Context,owner,repo,branchName,token string) error {	
	shaID,err := cuc.githubRepo.GetMasterSHA(owner,repo,token)
	if err != nil {
		return err
	}

	err = cuc.githubRepo.CreateBranch(owner,repo,branchName,shaID,token)
	if err != nil {
		return err
	}

	fileName, ok := ctx.Value(constants.FileNameKey).(string)
	if !ok {
		return errors.New("FileName not present")
	}

	fileContent,err := cuc.githubRepo.GetFileContent(owner,repo,branchName,fileName,token)
	if err != nil {
		return err
	}

	updateReq, err := GetUpdateFileReq(fileContent,branchName)
	if err!= nil {
		return err
	}

	err = cuc.githubRepo.CreateORUpdateFile(owner,repo,fileName,token, updateReq)
	if err!= nil {
		return err
	}

	req := GetPullReq(branchName, constants.MasterBranch,"Test Pull Request","Auto generated pull request")
	fmt.Println("Pull Req struct instance = ",req)
	err = cuc.githubRepo.CreatePullRequest(owner, repo, token , req)
	if err != nil {
		return err
	}

	return nil

}

func GetUpdateFileReq(fileContent *entity.FileContentResp, branchName string)  (*entity.UpdateFileReq,error) {
	updateReq := &entity.UpdateFileReq{}
	
	content := fileContent.Content
	lineCount := strings.Count(content,constants.NewLine) + 2
	line := fmt.Sprintf(constants.UpdateString,lineCount)

	updateReq.Message = fmt.Sprintf(constants.FileUpdateMsg,lineCount )
	updateReq.Content = base64.StdEncoding.EncodeToString([]byte(fileContent.Content + "\n" + line))
	updateReq.Branch = branchName
	updateReq.Sha =fileContent.Sha

	return updateReq,nil
}

func GetPullReq(head,base,body,title string) entity.PullReq {
	return entity.PullReq{
		Title: title,
		Head : head,
		Base : base,
		Body : body,
	}
}