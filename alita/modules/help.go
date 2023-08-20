package modules

import (
	"fmt"
	"html"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/divideprojects/Alita_Robot/alita/utils/helpers"
	log "github.com/sirupsen/logrus"

	"github.com/divideprojects/Alita_Robot/alita/config"
	"github.com/divideprojects/Alita_Robot/alita/db"
	"github.com/divideprojects/Alita_Robot/alita/i18n"

	"github.com/divideprojects/Alita_Robot/alita/utils/string_handling"
)

var HelpModule = moduleStruct{
	moduleName:     "Help",
	AbleMap:        moduleEnabled{},
	AltHelpOptions: make(map[string][]string),
	helpableKb:     make(map[string][][]gotgbot.InlineKeyboardButton),
}

const (
	aboutText = "Info & About\n\nHere are some of the FAQs about Alita."
	startHelp = "𝗛𝗲𝘆 𝗧𝗵𝗲𝗿𝗲! 𝗠𝘆 𝗡𝗮𝗺𝗲 𝗜𝘀 𝗕𝗲𝗲𝗿𝘂𝘀 😈\n" +
		"𝗜𝗺 𝗛𝗲𝗿𝗲 𝗧𝗼 𝗛𝗲𝗹𝗽 𝗬𝗼𝘂 𝗠𝗮𝗻𝗮𝗴𝗲 𝗬𝗼𝘂𝗿 𝗚𝗿𝗼𝘂𝗽!\n" +
		"𝗛𝗶𝘁 /mstart 𝗙𝗼𝗿 𝗠𝘂𝘀𝗶𝗰 𝗦𝗲𝗰𝘁𝗶𝗼𝗻\n" +
		"𝗝𝗼𝗶𝗻 𝗠𝘆 <a href='https://t.me/Godx_bots'>𝗡𝗲𝘄𝘀 𝗖𝗵𝗮𝗻𝗻𝗲𝗹</a> 𝗧𝗼 𝗚𝗲𝘁 𝗜𝗻𝗳𝗼𝗿𝗺𝗮𝘁𝗶𝗻𝗼 𝗢𝗻 𝗔𝗹𝗹 𝗧𝗵𝗲 𝗟𝗮𝘁𝗲𝘀𝘁 𝗨𝗽𝗱𝗮𝘁𝗲𝘀."
	mainhlp = "Hey %s!\n" +
		"My name is Alita ✨.\n\n" +
		"I am a group management bot, here to help you get around and keep the order in your groups!\n" +
		"I have lots of handy features, such as flood control, a warning system, a note keeping system, " +
		"and even predetermined replies on certain keywords.\n\n" +
		"<b>Helpful commands</b>:\n" +
		" - /start: Starts me! You've probably already used this!\n" +
		" - /help: Sends this message; I'll tell you more about myself!\n" +
		" - /donate: Gives you info on how to support me and my creator.\n\n" +
		"All commands can be used with the following: / or !"
)

var (
	backBtnSuffix = []gotgbot.InlineKeyboardButton{
		{
			Text:         "« Back",
			CallbackData: "helpq.Help",
		},
		{
			Text:         "Home",
			CallbackData: "helpq.BackStart",
		},
	}
	aboutKb = gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text:         "About me 👨\u200d💻",
					CallbackData: "about.me",
				},
			},
			{
				{
					Text: "News Channel 📢",
					Url:  "https://t.me/AlitaRobotUpdates",
				},
				{
					Text: "Support Group 👥",
					Url:  "https://t.me/DivideSupport",
				},
			},
			{
				{
					Text:         "Configuration ⚙️",
					CallbackData: "configuration.step1",
				},
			},
			{
				// custom back button
				{
					Text:         "⬅ Back",
					CallbackData: "helpq.BackStart",
				},
			},
		},
	}
	startMarkup = gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text:         "About ✨",
					CallbackData: "about.main",
				},
			},
			{
				{
					Text: "𝗔𝗱𝗱 𝗠𝗲",
					Url:  "https://t.me/beerus_godxbot?startgroup=botstart",
				},
				{
					Text: "𝗦𝘂𝗽𝗽𝗼𝗿𝘁",
					Url:  "https://t.me/the_strangers_home",
				},
			},
			{
				{
					Text:         "𝗛𝗲𝗹𝗽",
					CallbackData: "helpq.Help",
				},
			},
			{
				{
					Text:         "𝗟𝗮𝗻𝗴𝘂𝗮𝗴𝗲",
					CallbackData: "helpq.Languages",
				},
			},
		},
	}
)

type moduleEnabled struct {
	modules map[string]bool
}

func (m *moduleEnabled) Init() {
	m.modules = make(map[string]bool)
}

func (m *moduleEnabled) Store(module string, enabled bool) {
	m.modules[module] = enabled
}

func (m *moduleEnabled) Load(module string) (string, bool) {
	log.Info(fmt.Sprintf("[Module] Loading %s module", module))
	return module, m.modules[module]
}

func (m *moduleEnabled) LoadModules() []string {
	modules := make([]string, 0)
	for module := range m.modules {
		moduleName, enabled := m.Load(module)
		if enabled {
			modules = append(modules, moduleName)
		}
	}
	return modules
}

func (moduleStruct) about(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage

	tr := i18n.I18n{LangCode: db.GetLanguage(ctx)}

	var (
		currText string
		currKb   gotgbot.InlineKeyboardMarkup
	)

	if query := ctx.CallbackQuery; query != nil {
		args := strings.Split(query.Data, ".")
		response := args[1]

		switch response {
		case "main":
			currText = aboutText
			currKb = aboutKb
		case "me":
			currText = fmt.Sprintf(tr.GetString("strings.Help.About"), b.Username, config.BotVersion)
			currKb = gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{
							Text:         "Back",
							CallbackData: "about.main",
						},
					},
				},
			}
		}
		_, _, err := query.Message.EditText(b,
			currText,
			&gotgbot.EditMessageTextOpts{
				ReplyMarkup:           currKb,
				DisableWebPagePreview: true,
				ParseMode:             helpers.HTML,
			},
		)
		if err != nil {
			log.Error(err)
			return err
		}

		_, err = query.Answer(b, nil)
		if err != nil {
			log.Error(err)
			return err
		}
	} else {
		if ctx.Update.Message.Chat.Type == "private" {
			currText = aboutText
			currKb = aboutKb
		} else {
			currText = "Click on the button below to get info about me!"
			currKb = gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{
							Text: "About",
							Url:  fmt.Sprintf("https://t.me/%s?start=about", b.Username),
						},
					},
				},
			}
		}
		_, err := msg.Reply(
			b,
			currText,
			&gotgbot.SendMessageOpts{
				ParseMode:             helpers.HTML,
				DisableWebPagePreview: true,
				ReplyMarkup:           &currKb,
			},
		)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	return ext.EndGroups
}

func (moduleStruct) helpButtonHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.Update.CallbackQuery
	args := strings.Split(query.Data, ".")
	module := args[1]

	var (
		_parsemode, helpText string
		replyKb              gotgbot.InlineKeyboardMarkup
	)

	// Sort the module names
	if string_handling.FindInStringSlice([]string{"BackStart", "Help"}, module) {
		_parsemode = helpers.HTML
		switch module {
		case "Help":
			// This shows the main start menu
			helpText = fmt.Sprintf(mainhlp, html.EscapeString(query.From.FirstName))
			replyKb = markup
		case "BackStart":
			// This shows the modules menu
			helpText = startHelp
			replyKb = startMarkup
		}
	} else {
		// For all remainging modules
		helpText, replyKb, _parsemode = getHelpTextAndMarkup(ctx, strings.ToLower(module))
	}

	// Edit the main message, the main querymessage
	_, _, err := query.Message.EditText(
		b,
		helpText,
		&gotgbot.EditMessageTextOpts{
			ParseMode:             _parsemode,
			ReplyMarkup:           replyKb,
			DisableWebPagePreview: true,
		},
	)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = query.Answer(b, nil)
	if err != nil {
		log.Error(err)
		return err
	}

	return ext.EndGroups
}

// start introduces the bot
func (moduleStruct) start(b *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveSender.User
	msg := ctx.EffectiveMessage
	args := ctx.Args()

	if ctx.Update.Message.Chat.Type == "private" {
		if len(args) == 1 {
			_, err := msg.Reply(b,
				startHelp,
				&gotgbot.SendMessageOpts{
					ParseMode:             helpers.HTML,
					DisableWebPagePreview: true,
					ReplyMarkup:           &startMarkup,
				},
			)
			if err != nil {
				log.Error(err)
				return err
			}
		} else if len(args) == 2 {
			err := startHelpPrefixHandler(b, ctx, user, args[1])
			if err != nil {
				log.Error(err)
				return err
			}
		} else {
			log.Info("sed")
		}
	} else {
		_, err := msg.Reply(b, "Hey :) PM me if you have any questions on how to use me!", helpers.Shtml())
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return ext.EndGroups
}

func (moduleStruct) donate(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	chat := ctx.EffectiveChat

	_, err := b.SendMessage(chat.Id,
		i18n.I18n{LangCode: "en"}.GetString("strings.Help.DonateText"),
		&gotgbot.SendMessageOpts{
			ParseMode:                helpers.HTML,
			DisableWebPagePreview:    true,
			ReplyToMessageId:         msg.MessageId,
			AllowSendingWithoutReply: true,
		},
	)
	if err != nil {
		log.Error(err)
	}

	return ext.EndGroups
}

func (moduleStruct) botConfig(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.CallbackQuery
	msg := query.Message

	// just in case
	if msg.Chat.Type != "private" {
		_, _, err := msg.EditText(b, "Configuration only works in private", nil)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	args := strings.Split(query.Data, ".")
	response := args[1]

	var (
		iKeyboard [][]gotgbot.InlineKeyboardButton
		text      string
	)

	tr := i18n.I18n{LangCode: "en"}

	switch response {
	case "step1":
		iKeyboard = [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text: "➕ Add Alita to chat!",
					Url:  fmt.Sprintf("https://t.me/%s?startgroup=botstart", b.Username),
				},
			},
			{
				{
					Text:         "Done ✅",
					CallbackData: "configuration.step2",
				},
			},
		}
		text = tr.GetString("strings.Help.Configuration.Step-1")
	case "step2":
		iKeyboard = [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text:         "Done ✅",
					CallbackData: "configuration.step3",
				},
			},
		}
		text = fmt.Sprintf(tr.GetString("strings.Help.Configuration.Step-2"), b.Username)
	case "step3":
		iKeyboard = [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text:         "✅ Continue ✅",
					CallbackData: "helpq.Help",
				},
			},
		}
		text = tr.GetString("strings.Help.Configuration.Step-3")
	}
	_, _, err := msg.EditText(
		b,
		text,
		&gotgbot.EditMessageTextOpts{
			DisableWebPagePreview: true,
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: iKeyboard,
			},
		},
	)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = query.Answer(b, nil)
	if err != nil {
		log.Error(err)
		return err
	}

	return ext.EndGroups
}

func (moduleStruct) help(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	msg := ctx.EffectiveMessage
	args := ctx.Args()

	if ctx.Update.Message.Chat.Type == "private" {
		if len(args) == 1 {
			_, err := b.SendMessage(chat.Id,
				fmt.Sprintf(
					mainhlp,
					html.EscapeString(msg.From.FirstName),
				),
				&gotgbot.SendMessageOpts{
					ParseMode:   helpers.HTML,
					ReplyMarkup: &markup,
				},
			)
			if err != nil {
				log.Error(err)
				return err
			}
		} else if len(args) == 2 {
			module := strings.ToLower(args[1])
			helpText, replyMarkup, _parsemode := getHelpTextAndMarkup(ctx, module)
			_, err := b.SendMessage(
				chat.Id,
				helpText,
				&gotgbot.SendMessageOpts{
					ParseMode:   _parsemode,
					ReplyMarkup: &replyMarkup,
				},
			)
			if err != nil {
				log.Error(err)
				return err
			}
		}
	} else {
		pmMeKbText := "Click here for help!"
		pmMeKbUri := fmt.Sprintf("https://t.me/%s?start=help_help", b.Username)
		moduleHelpString := "Contact me in PM for help!"
		replyMsgId := msg.MessageId
		var lowerModName string

		if len(args) == 2 {
			helpModName := args[1]
			lowerModName = strings.ToLower(helpModName)
			originalModuleName := getModuleNameFromAltName(lowerModName)
			if originalModuleName != "" && string_handling.FindInStringSlice(getAltNamesOfModule(originalModuleName), lowerModName) {
				moduleHelpString = fmt.Sprintf("Contact me in PM for help regarding <code>%s</code>!", originalModuleName)
				pmMeKbUri = fmt.Sprintf("https://t.me/%s?start=help_%s", b.Username, lowerModName)
			}
		}

		if msg.ReplyToMessage != nil {
			replyMsgId = msg.ReplyToMessage.MessageId
		}

		_, err := msg.Reply(b,
			moduleHelpString,
			&gotgbot.SendMessageOpts{
				ParseMode:                helpers.HTML,
				ReplyToMessageId:         replyMsgId,
				AllowSendingWithoutReply: true,
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{
					InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
						{
							{Text: pmMeKbText, Url: pmMeKbUri},
						},
					},
				},
			},
		)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return ext.EndGroups
}

func LoadHelp(dispatcher *ext.Dispatcher) {
	dispatcher.AddHandler(handlers.NewCommand("start", HelpModule.start))
	dispatcher.AddHandler(handlers.NewCommand("help", HelpModule.help))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("helpq"), HelpModule.helpButtonHandler))
	dispatcher.AddHandler(handlers.NewCommand("donate", HelpModule.donate))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("configuration"), HelpModule.botConfig))
	dispatcher.AddHandler(handlers.NewCommand("about", HelpModule.about))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("about"), HelpModule.about))
	initHelpButtons()
}
