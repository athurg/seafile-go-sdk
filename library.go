package seafile

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	LibraryTypeMine   = "mine"   //我的资料库类型
	LibraryTypeShared = "shared" //私人共享给我的资料库类型
	LibraryTypeGroup  = "group"  //群组共享给我的资料库类型
	LibraryTypeOrg    = "org"    //公共资料库类型
)

//资料库
type Library struct {
	Id         string
	Name       string
	Type       string
	Root       string
	Owner      string
	Permission string

	Encrypted bool
	Virtual   bool
	Version   int
	Mtime     int
	Size      int

	MtimeRelative string `json:"mtime_relative"`
	HeadCommitId  string `json:"head_commit_id"`
	SizeFormatted string `json:"size_formatted"`

	client *Client
}

//获取可用的资料库
func (cli *Client) ListAllLibraries() ([]*Library, error) {
	return cli.ListLibrariesByType("")
}

//获取拥有的资料库
func (cli *Client) ListOwnedLibraries() ([]*Library, error) {
	return cli.ListLibrariesByType(LibraryTypeMine)
}

//获取私人共享而来的资料库
func (cli *Client) ListSharedLibraries() ([]*Library, error) {
	return cli.ListLibrariesByType(LibraryTypeShared)
}

//获取群组共享而来的资料库
func (cli *Client) ListGroupLibraries() ([]*Library, error) {
	return cli.ListLibrariesByType(LibraryTypeGroup)
}

//获取公共的资料库
func (cli *Client) ListOrgLibraries() ([]*Library, error) {
	return cli.ListLibrariesByType(LibraryTypeOrg)
}

//获取指定类型的资料库
func (cli *Client) ListLibrariesByType(libType string) ([]*Library, error) {
	uri := "/repos/"
	if libType != "" {
		uri += "?type=" + libType
	}

	resp, err := cli.doRequest("GET", uri, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	info := []*Library{}
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil, fmt.Errorf("读取错误:%s %s", resp.Status, err)
	}

	for _, lib := range info {
		lib.client = cli
	}

	return info, nil
}

//获取默认资料库
func (cli *Client) GetDefaultLibrary() (*Library, error) {
	return cli.GetLibrary("")
}

func (cli *Client) GetLibrary(name string) (*Library, error) {
	var id string
	var err error
	//如果name为空字符串，则获取默认资料库
	if name == "" {
		id, err = cli.GetDefaultLibraryId()
		if err != nil {
			return nil, err
		}
	}

	libraries, err := cli.ListAllLibraries()
	if err != nil {
		return nil, fmt.Errorf("获取资料库列表失败: %s", err)
	}

	if name == "" {
		for _, library := range libraries {
			if library.Id == id {
				return library, nil
			}
		}
	} else {
		for _, library := range libraries {
			if library.Name == name {
				return library, nil
			}
		}
	}

	return nil, fmt.Errorf("未找到资料库")
}

//获取默认资料库ID
func (cli *Client) GetDefaultLibraryId() (string, error) {
	resp, err := cli.doRequest("GET", "/default-repo/", nil, nil)
	if err != nil {
		return "", fmt.Errorf("获取默认资料库失败: %s", err)
	}
	defer resp.Body.Close()

	var respInfo struct {
		Exists bool
		RepoId string `json:"repo_id"`
	}

	err = json.NewDecoder(resp.Body).Decode(&respInfo)
	if err != nil {
		return "", fmt.Errorf("获取默认资料库失败: %s", err)
	}

	if !respInfo.Exists {
		return "", fmt.Errorf("默认资料库不存在")
	}

	return respInfo.RepoId, nil
}

func (lib *Library) doRequest(method, uri string, header http.Header, body io.Reader) (*http.Response, error) {
	if !strings.HasPrefix(uri, "http://") && !strings.HasPrefix(uri, "https://") {
		uri = "/repos/" + lib.Id + uri
	}
	return lib.client.doRequest(method, uri, header, body)
}

//获取资料库的上传地址
func (lib *Library) UploadLink() (string, error) {
	resp, err := lib.doRequest("GET", "/upload-link/", nil, nil)
	if err != nil {
		return "", fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	var link string
	err = json.NewDecoder(resp.Body).Decode(&link)
	if err != nil {
		return "", fmt.Errorf("解析错误:%s %s", resp.Status, err)
	}

	//返回值是"xxx"格式的，需要去掉头尾的引号
	return link, nil
}

//资料库提交
type LibraryCommit struct {
	Id                string
	Desc              string
	Ctime             int
	Creator           string
	Conflict          bool
	NewMerge          bool   `json:"new_merge"`
	CreatorName       string `json:"creator_name"`
	RootId            string `json:"root_id"`
	RepoId            string `json:"repo_id"`
	ParentId          string `json:"parent_id"`
	SecondParentId    string `json:"second_parent_id"`
	RevFileSize       int    `json:"rev_file_size"`
	RevFileId         string `json:"rev_file_id"`
	RevRenamedOldPath string `json:"rev_renamed_old_path"`

	library *Library `json:"-"`
}

//获取资料库的提交历史
func (lib *Library) History() ([]*LibraryCommit, error) {
	resp, err := lib.doRequest("GET", "/history", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取错误: %s", err)
	}

	var respInfo struct {
		PageNext bool `json:"page_next"`
		Commits  []*LibraryCommit
	}

	err = json.Unmarshal(b, &respInfo)
	if err != nil {
		return nil, fmt.Errorf("解析错误: %s %s", err, string(b))
	}

	for _, commit := range respInfo.Commits {
		commit.library = lib
	}

	return respInfo.Commits, nil
}
