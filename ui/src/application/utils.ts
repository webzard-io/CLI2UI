import { ComponentSchema } from "@sunmao-ui/core";

const firstUpperCase = (str: string) => {
  return str.toLowerCase().replace(/( |^)[a-z]/g, (L) => L.toUpperCase());
};

enum ArgsType {
  String = "string",
  Number = "number",
}

type CLIJson = {
  title: string;
  help: string;
  cmd: {
    name: string;
    args: {
      name: string;
      type: ArgsType;
      required: boolean;
    }[];
  }[];
};

const raw: CLIJson = {
  title: "ping",
  help: "usage: ping [-c count] host",
  cmd: [
    {
      name: "ping",
      args: [
        {
          name: "count",
          required: false,
          type: ArgsType.Number,
        },
        {
          name: "host",
          required: true,
          type: ArgsType.String,
        },
      ],
    },
    {
      name: "ping2",
      args: [],
    },
  ],
};

export const genSchemaComponents = () => {
  const genLayout = () => {
    return [
      {
        id: "Layout",
        type: "arco/v1/layout",
        properties: {
          showHeader: true,
          showSideBar: false,
          sidebarCollapsible: false,
          sidebarDefaultCollapsed: false,
          showFooter: false,
        },
        traits: [
          {
            type: "core/v1/style",
            properties: {
              styles: [
                {
                  styleSlot: "layout",
                  style: "\n",
                  cssProperties: {
                    paddingLeft: "32px",
                    paddingRight: "32px",
                    paddingTop: "16px",
                    paddingBottom: "16px",
                    marginRight: "",
                    backgroundColor: "",
                    marginBottom: "",
                  },
                },
              ],
            },
          },
        ],
      },
    ];
  };

  const genHeader = () => {
    return [
      {
        id: "Header",
        type: "core/v1/stack",
        properties: {
          spacing: 12,
          direction: "horizontal",
          align: "center",
          wrap: false,
          justify: "flex-start",
        },
        traits: [
          {
            type: "core/v1/style",
            properties: {
              styles: [
                {
                  styleSlot: "content",
                  style: "",
                  cssProperties: {
                    width: "100%",
                  },
                },
              ],
            },
          },
          {
            type: "core/v2/slot",
            properties: {
              container: {
                id: "Layout",
                slot: "header",
              },
              ifCondition: true,
            },
          },
        ],
      },
      {
        id: "HeaderText",
        type: "core/v2/text",
        properties: {
          text: raw.title,
        },
        traits: [
          {
            type: "core/v2/slot",
            properties: {
              container: {
                id: "Header",
                slot: "content",
              },
              ifCondition: true,
            },
          },
          {
            type: "core/v1/style",
            properties: {
              styles: [
                {
                  styleSlot: "content",
                  style: "",
                  cssProperties: {
                    fontSize: "32px\n",
                    fontWeight: "700",
                  },
                },
              ],
            },
          },
        ],
      },
      {
        id: "HelpBtn",
        type: "arco/v1/button",
        properties: {
          type: "default",
          status: "default",
          long: false,
          size: "default",
          disabled: false,
          loading: false,
          shape: "square",
          text: "查看",
        },
        traits: [
          {
            type: "core/v2/slot",
            properties: {
              container: {
                id: "Header",
                slot: "content",
              },
              ifCondition: true,
            },
          },
          {
            type: "core/v1/event",
            properties: {
              handlers: [
                {
                  type: "onClick",
                  componentId: "HelpModal",
                  method: {
                    name: "openModal",
                    parameters: {},
                  },
                  wait: {
                    type: "debounce",
                    time: 0,
                  },
                  disabled: false,
                },
              ],
            },
          },
        ],
      },
    ];
  };

  const genHelpModal = () => {
    return [
      {
        id: "HelpModal",
        type: "arco/v1/modal",
        properties: {
          title: "Help 信息",
          mask: true,
          simple: false,
          okText: "confirm",
          cancelText: "cancel",
          closable: true,
          maskClosable: true,
          confirmLoading: false,
          defaultOpen: false,
          unmountOnExit: true,
        },
        traits: [],
      },
      {
        id: "HelpModalCancelBtn",
        type: "arco/v1/button",
        properties: {
          type: "default",
          status: "default",
          long: false,
          size: "default",
          disabled: false,
          loading: false,
          shape: "square",
          text: "关闭",
        },
        traits: [
          {
            type: "core/v2/slot",
            properties: {
              container: {
                id: "HelpModal",
                slot: "footer",
              },
              ifCondition: true,
            },
          },
          {
            type: "core/v1/event",
            properties: {
              handlers: [
                {
                  type: "onClick",
                  componentId: "HelpModal",
                  method: {
                    name: "closeModal",
                    parameters: {},
                  },
                  wait: {
                    type: "debounce",
                    time: 0,
                  },
                  disabled: false,
                },
              ],
            },
          },
          {
            type: "core/v1/style",
            properties: {
              styles: [],
            },
          },
        ],
      },
      {
        id: "HelpModalContent",
        type: "core/v1/stack",
        properties: {
          spacing: 12,
          direction: "horizontal",
          align: "auto",
          wrap: false,
          justify: "flex-start",
        },
        traits: [
          {
            type: "core/v2/slot",
            properties: {
              container: {
                id: "HelpModal",
                slot: "content",
              },
              ifCondition: true,
            },
          },
          {
            type: "core/v1/style",
            properties: {
              styles: [
                {
                  styleSlot: "content",
                  style: "",
                  cssProperties: {
                    backgroundColor: "#333",
                    color: "white",
                    width: "100%",
                    paddingTop: "16px",
                    paddingLeft: "16px",
                    paddingBottom: "16px",
                    paddingRight: "16px",
                  },
                },
              ],
            },
          },
        ],
      },
      {
        id: "HelpInfo",
        type: "core/v2/text",
        properties: {
          text: `{{"${raw.help}"}}`,
        },
        traits: [
          {
            type: "core/v2/slot",
            properties: {
              container: {
                id: "HelpModalContent",
                slot: "content",
              },
              ifCondition: true,
            },
          },
          {
            type: "core/v1/style",
            properties: {
              styles: [],
            },
          },
        ],
      },
    ];
  };

  const genCmdTab = () => {
    const tab = {
      id: "CmdTabs",
      type: "arco/v1/tabs",
      properties: {
        type: "line",
        defaultActiveTab: 0,
        tabPosition: "top",
        size: "default",
        updateWhenDefaultValueChanges: false,
        tabs: raw.cmd.map((item) => ({
          title: item.name,
          hidden: false,
          destroyOnHide: true,
        })),
      },
      traits: [
        {
          type: "core/v2/slot",
          properties: {
            container: {
              id: "Layout",
              slot: "content",
            },
            ifCondition: true,
          },
        },
      ],
    };

    return [tab];
  };

  const genCmdFormFiledInput = (
    cmdName: string,
    arg: CLIJson["cmd"][0]["args"][0]
  ) => {
    switch (arg.type) {
      case ArgsType.Number:
        return {
          id: `${cmdName}Form${firstUpperCase(arg.name)}Input`,
          type: "arco/v1/numberInput",
          properties: {
            defaultValue: 1,
            disabled: `{{${cmdName}FormState.data.isExecuting}}`,
            placeholder: "please input",
            error: false,
            size: "default",
            buttonMode: false,
            min: 0,
            max: 99,
            readOnly: false,
            step: 1,
            precision: 0,
            updateWhenDefaultValueChanges: false,
          },
          traits: [
            {
              type: "core/v2/slot",
              properties: {
                container: {
                  id: `${cmdName}Form${firstUpperCase(arg.name)}Field`,
                  slot: "content",
                },
                ifCondition: true,
              },
            },
            {
              type: "core/v1/event",
              properties: {
                handlers: [],
              },
            },
          ],
        };
      default:
        return {
          id: `${cmdName}Form${firstUpperCase(arg.name)}Input`,
          type: "arco/v1/input",
          properties: {
            allowClear: false,
            disabled: `{{${cmdName}FormState.data.isExecuting}}`,
            readOnly: false,
            defaultValue: "",
            updateWhenDefaultValueChanges: false,
            placeholder: "please input",
            error: `{{${cmdName}Form.validatedResult.${cmdName}Form${firstUpperCase(
              arg.name
            )}Validation.isInvalid}}`,
            size: "default",
          },
          traits: [
            {
              type: "core/v2/slot",
              properties: {
                container: {
                  id: `${cmdName}Form${firstUpperCase(arg.name)}Field`,
                  slot: "content",
                },
                ifCondition: true,
              },
            },
            {
              type: "core/v1/event",
              properties: {
                handlers: [
                  {
                    type: "onBlur",
                    componentId: `${cmdName}Form`,
                    method: {
                      name: "validateAllFields",
                      parameters: {},
                    },
                    wait: {
                      type: "debounce",
                      time: 0,
                    },
                    disabled: false,
                  },
                ],
              },
            },
            {
              type: "core/v1/style",
              properties: {
                styles: [
                  {
                    styleSlot: "input",
                    style: "",
                    cssProperties: {
                      marginBottom: "6px",
                    },
                  },
                ],
              },
            },
          ],
        };
    }
  };

  const genCmdFormFields = (cmdName: string, args: CLIJson["cmd"][0]["args"]) => {
    let formFieldComponents = [] as unknown[];
    args.forEach((arg) => {
      const components = [
        {
          id: `${cmdName}Form${firstUpperCase(arg.name)}Field`,
          type: "arco/v1/formControl",
          properties: {
            label: {
              format: "plain",
              raw: `${arg.name}`,
            },
            layout: "horizontal",
            required: false,
            hidden: false,
            extra: "",
            errorMsg: "",
            labelAlign: "left",
            colon: false,
            help: "",
            labelCol: {
              span: 3,
              offset: 0,
            },
            wrapperCol: {
              span: 21,
              offset: 0,
            },
          },
          traits: [
            {
              type: "core/v2/slot",
              properties: {
                container: {
                  id: `${cmdName}Form`,
                  slot: "content",
                },
                ifCondition: true,
              },
            },
          ],
        },
        genCmdFormFiledInput(cmdName, arg),
      ];
      formFieldComponents = formFieldComponents.concat(components);
    });
    return formFieldComponents;
  };

  const genCmdInnerTabs = () => {
    let tabs = [] as unknown[];

    raw.cmd.forEach((item, index) => {
      const tabComponents = [
        {
          id: `${item.name}FormState`,
          type: "core/v1/dummy",
          properties: {},
          traits: [
            {
              type: "core/v1/state",
              properties: {
                key: "data",
                initialValue: "{{\n{\n  isExecuting: false,\n}\n}}",
              },
            },
            {
              type: "core/v1/event",
              properties: {
                handlers: [],
              },
            },
          ],
        },
        {
          id: `${item.name}Tab`,
          type: "core/v1/stack",
          properties: {
            spacing: 12,
            direction: "vertical",
            align: "auto",
            wrap: false,
            justify: "flex-start",
          },
          traits: [
            {
              type: "core/v2/slot",
              properties: {
                container: {
                  id: "CmdTabs",
                  slot: "content",
                },
                ifCondition: `{{CmdTabs.activeTab === ${index}}}`,
              },
            },
            {
              type: "core/v1/style",
              properties: {
                styles: [
                  {
                    styleSlot: "content",
                    style: "",
                    cssProperties: {
                      width: "100%",
                    },
                  },
                ],
              },
            },
          ],
        },
        {
          id: `${item.name}Form`,
          type: "core/v1/stack",
          properties: {
            spacing: 12,
            direction: "vertical",
            align: "auto",
            wrap: false,
            justify: "flex-start",
          },
          traits: [
            {
              type: "core/v2/slot",
              properties: {
                container: {
                  id: `${item.name}Tab`,
                  slot: "content",
                },
                ifCondition: true,
              },
            },
            {
              type: "core/v1/validation",
              properties: {
                validators: item.args.map((sub) => {
                  const validation = {
                    name: `${item.name}Form${firstUpperCase(sub.name)}Validation`,
                    value: `{{${item.name}Form${firstUpperCase(sub.name)}Input.value}}`,
                    rules: [] as unknown[],
                  };
                  if (sub.required) {
                    validation.rules.push({
                      type: "required",
                      validate: "",
                      error: {
                        message: "该字段不可为空。",
                      },
                      minLength: 0,
                      maxLength: 0,
                      includeList: [],
                      excludeList: [],
                      min: 0,
                      max: 0,
                      regex: "",
                      flags: "",
                      customOptions: {},
                    });
                  }
                  return validation;
                }),
              },
            },
          ],
        },
        ...genCmdFormFields(item.name, item.args),
        {
          id: `${item.name}Button`,
          type: "arco/v1/button",
          properties: {
            type: "primary",
            status: "default",
            long: false,
            size: "default",
            disabled: false,
            loading: `{{${item.name}FormState.data.isExecuting}}`,
            shape: "square",
            text: "执行",
          },
          traits: [
            {
              type: "core/v2/slot",
              properties: {
                container: {
                  id: `${item.name}Form`,
                  slot: "content",
                },
                ifCondition: true,
              },
            },
            {
              type: "core/v1/event",
              properties: {
                handlers: [
                  {
                    type: "onClick",
                    componentId: `${item.name}Form`,
                    method: {
                      name: "validateAllFields",
                      parameters: {},
                    },
                    wait: {
                      type: "debounce",
                      time: 0,
                    },
                    disabled: false,
                  },
                  {
                    type: "onClick",
                    componentId: `${item.name}FormState`,
                    method: {
                      name: "setValue",
                      parameters: {
                        key: "data",
                        value: `{{{...${item.name}FormState.data, isExecuting: true}}}`,
                      },
                    },
                    wait: {
                      type: "debounce",
                      time: 0,
                    },
                    disabled: `{{${item.name}Form.isInvalid}}`,
                  },
                ],
              },
            },
          ],
        },
        {
          id: `${item.name}TabDivider`,
          type: "arco/v1/divider",
          properties: {
            type: "horizontal",
            orientation: "center",
          },
          traits: [
            {
              type: "core/v2/slot",
              properties: {
                container: {
                  id: `${item.name}Tab`,
                  slot: "content",
                },
                ifCondition: true,
              },
            },
          ],
        },
        {
          id: `${item.name}TabResult`,
          type: "arco/v1/collapse",
          properties: {
            defaultActiveKey: '{{[\n  "0"\n]}}',
            options: [
              {
                key: "0",
                header: "执行结果\n",
                disabled: false,
                showExpandIcon: true,
              },
            ],
            updateWhenDefaultValueChanges: false,
            accordion: false,
            expandIconPosition: "left",
            bordered: false,
            destroyOnHide: false,
            lazyLoad: true,
          },
          traits: [
            {
              type: "core/v2/slot",
              properties: {
                container: {
                  id: `${item.name}Tab`,
                  slot: "content",
                },
                ifCondition: true,
              },
            },
            {
              type: "core/v1/style",
              properties: {
                styles: [
                  {
                    styleSlot: "content",
                    style:
                      ".arco-collapse-item-content-box {\nbackground: white;\n}",
                    cssProperties: {
                      backgroundColor: "",
                    },
                  },
                ],
              },
            },
          ],
        },
        {
          id: `${item.name}StopBtn`,
          type: "arco/v1/button",
          properties: {
            type: "primary",
            status: "default",
            long: false,
            size: "default",
            disabled: `{{!${item.name}FormState.data.isExecuting}}`,
            loading: false,
            shape: "square",
            text: "停止",
          },
          traits: [
            {
              type: "core/v2/slot",
              properties: {
                container: {
                  id: `${item.name}Tab`,
                  slot: "content",
                },
                ifCondition: true,
              },
            },
            {
              type: "core/v1/event",
              properties: {
                handlers: [
                  {
                    type: "onClick",
                    componentId: `${item.name}FormState`,
                    method: {
                      name: "setValue",
                      parameters: {
                        key: "data",
                        value: `{{{...${item.name}FormState.data, isExecuting: false}}}`,
                      },
                    },
                    wait: {
                      type: "debounce",
                      time: 0,
                    },
                    disabled: false,
                  },
                ],
              },
            },
            {
              type: "core/v1/style",
              properties: {
                styles: [
                  {
                    styleSlot: "content",
                    style: "",
                    cssProperties: {
                      width: "",
                    },
                  },
                ],
              },
            },
          ],
        },
      ];
      tabs = tabs.concat(tabComponents);
    });

    return tabs;
  };

  return [
    ...genLayout(),
    ...genHeader(),
    ...genCmdTab(),
    ...genCmdInnerTabs(),
    ...genHelpModal(),
  ] as ComponentSchema[];
};
