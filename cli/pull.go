/*
File: pull.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-03-05 14:22:54

Description: 子命令 `pull` 的实现
*/

package cli

import (
	"fmt"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/yhyj/curator/general"
)

// RollingPullRepos 遍历拉取远端仓库的更改到本地
//
// 参数：
//   - confile: 程序配置文件
//   - source: 远端仓库源，支持 'github' 和 'gitea'，默认为 'github'
func RollingPullRepos(confile, source string) {
	// 加载配置文件
	conf, err := GetTomlConfig(confile)
	if err != nil {
		fmt.Printf(general.ErrorBaseFormat, err)
	} else {
		// 获取配置项
		pemfile := conf.Get("ssh.rsa_file")
		storagePath := conf.Get("storage.path").(string)
		repoNames := conf.Get("git.repos").([]interface{})
		// 获取公钥
		publicKeys, err := general.GetPublicKeysByGit(pemfile.(string))
		if err != nil {
			fmt.Printf(general.ErrorBaseFormat, err)
			return
		}

		// 拉取
		fmt.Printf(general.TipsPrefixSuffixFormat, "Pull changes from", " ", source, " ", "remote repository")
		fmt.Println()
		for _, repoName := range repoNames {
			repoPath := filepath.Join(storagePath, repoName.(string))
			// 拉取前检测本地仓库是否存在
			if general.FileExist(repoPath) {
				isRepo, repo := general.IsLocalRepo(repoPath)
				if isRepo { // 开始拉取
					fmt.Printf(general.Tips2PSuffixNoNewLineFormat, general.Run, " Pulling ", repoName.(string), ":", " ")
					leftCommit, rightCommit, err := general.PullRepo(repo, publicKeys)
					if err != nil {
						if err == git.NoErrAlreadyUpToDate {
							fmt.Println("Already up-to-date")
						} else {
							fmt.Printf(general.ErrorBaseFormat, err)
						}
					} else {
						fmt.Printf(general.SliceTraverse2PFormat, leftCommit.Hash.String()[:6], " --> ", rightCommit.Hash.String()[:6])
					}
				} else { // 非本地仓库
					fmt.Printf(general.Tips2PSuffixNoNewLineFormat, general.No, " Pulling ", repoName.(string), ":", " ")
					fmt.Println("Folder is not a local repository")
				}
			} else {
				fmt.Printf(general.Tips2PSuffixNoNewLineFormat, general.No, " Pulling ", repoName.(string), ":", " ")
				fmt.Println("The local repository does not exist")
			}
			// 添加一个延时，使输出更加顺畅
			general.Delay(0.1)
			continue
		}
	}
}
