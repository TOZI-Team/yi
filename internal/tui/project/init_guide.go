package project

import (
	"fmt"
	"github.com/erikgeiser/promptkit/textinput"
	log "github.com/sirupsen/logrus"
	"yi/internal/sdk"
	"yi/internal/tui/box/compiler"
	cjpackage "yi/pkg/package"
	t "yi/pkg/types"
)

func InitGuide(c t.InitConfig) t.InitConfig {

	// 询问包名
	packageInput := textinput.New(":: 请输入包名")
	packageInput.InitialValue = c.Name
	packageInput.Validate = func(input string) error {
		if cjpackage.IsCjPackageName(input) {
			return nil
		}
		return fmt.Errorf("非法包名")
	}
	packageName, err := packageInput.RunPrompt()
	if err != nil {
		log.Fatal(err)
	}
	c.Name = packageName

	// 询问包描述
	descriptionInput := textinput.New(":: 请输入包的描述：")
	descriptionInput.InitialValue = c.Description
	descriptionInput.Placeholder = "不能为空"
	desName, err := descriptionInput.RunPrompt()
	if err != nil {
		log.Fatal(err)
	}
	c.Description = desName

	// 询问版本
	verInput := textinput.New(":: 请输入包版本：")
	verInput.InitialValue = c.Version
	ver, err := verInput.RunPrompt()
	if err != nil {
		log.Fatal(err)
	}
	c.Version = ver

	com := compiler.ChooseCompiler(sdk.GlobalSDKManger.GetSDKs())
	c.ComVer = com.Ver
	s := new(t.SDKInfo)
	*(s) = com
	c.SDK = s
	c.Output = t.EXECUTABLE
	log.Infof("Set compiler vsrsion: %s", c.ComVer)

	// 初始化
	//ch := make(chan t.WaitingMessage, 1)
	//p := tea.NewProgram(modules.NewWaitingModel("初始化中 ...", ch))
	//if _, err := p.Run(); err != nil {
	//	log.Fatal(err)
	//}

	return c
}
