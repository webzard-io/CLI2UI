import { ComponentSchema, Application } from "@sunmao-ui/core";
import camelCase from "camelcase";
import kebabCase from "kebab-case";

const kebabToPascalCase = (str: string) => {
  return camelCase(str, { pascalCase: true });
};

const pascalCaseToKebabCase = (str: string) => {
  return kebabCase(str);
};

const componentIdReg =
  /^CLI2UI_QWE_([a-zA-Z0-9]+)_ASD_.*_QWE_([a-zA-Z0-9]+)_ASD_/;

const getFlagValidationId = (cmdName: string, flagName: string) => {
  return `CLI2UI_QEW_${cmdName}_ASD_Form_QEW_${kebabToPascalCase(
    flagName
  )}_ASD_Validation`;
};
const getFlagFieldId = (cmdName: string, flagName: string) => {
  return `CLI2UI_QEW_${cmdName}_ASD_Form_QEW_${kebabToPascalCase(
    flagName
  )}_ASD_Field`;
};
const getFlagInputId = (cmdName: string, flagName: string) => {
  return `CLI2UI_QEW_${cmdName}_ASD_Form_QEW_${kebabToPascalCase(
    flagName
  )}_ASD_Input`;
};

export enum FlagType {
  String = "string",
  Number = "number",
  Array = "array",
  Boolean = "boolean",
  Enum = "enum",
}

export type CLIJson = {
  title: string;
  help: string;
  cmd: {
    name: string;
    flags: {
      name: string;
      type: FlagType;
      required?: boolean;
      options?: string[];
    }[];
  }[];
};

const pingRaw: CLIJson = {
  title: "ping",
  help: `usage: ping [-AaDdfnoQqRrv] [-c count] [-G sweepmaxsize]
                [-g sweepminsize] [-h sweepincrsize] [-i wait]
                [-l preload] [-M mask | time] [-m ttl] [-p pattern]
                [-S src_addr] [-s packetsize] [-t timeout][-W waittime]
                [-z tos] host`,
  cmd: [
    {
      name: "ping",
      flags: [
        {
          name: "count",
          required: false,
          type: FlagType.Number,
        },
        {
          name: "host",
          required: true,
          type: FlagType.String,
        },
      ],
    },
  ],
};

const kailRaw: CLIJson = {
  title: "kail",
  help: `usage: kail [<flags>] <command> [<args> ...]
  
    Tail for kubernetes pods
    
    Flags:
      -h, --help                  Show context-sensitive help (also try --help-long
                                  and --help-man).
          --ignore=SELECTOR ...   ignore selector
      -l, --label=SELECTOR ...    label
      -p, --pod=NAME ...          pod
      -n, --ns=NAME ...           namespace
          --ignore-ns=NAME ...    ignore namespace
          --svc=NAME ...          service
          --rc=NAME ...           replication controller
          --rs=NAME ...           replica set
          --ds=NAME ...           daemonset
      -d, --deploy=NAME ...       deployment
          --sts=NAME ...          statefulset
      -j, --job=NAME ...          job
          --node=NAME ...         node
          --ing=NAME ...          ingress
          --context=CONTEXT-NAME  kubernetes context
          --current-ns            use namespace from current context
      -c, --containers=NAME ...   containers
          --dry-run               print matching pods and exit
          --log-file=LOG-FILE     log file output
          --log-level=error       log level
          --since=DURATION        Display logs generated since given duration,
                                  like 5s, 2m, 1.5h or 2h45m. Defaults to 1s.
      -o, --output=default        Log output mode (default, raw, json, or
                                  json-pretty, zerolog)
          --zerolog-timestamp-field="time"
                                  sets the zerolog timestamp field name, works with
                                  --output=zerolog
          --zerolog-level-field="level"
                                  sets the zerolog level field name, works with
                                  --output=zerolog
          --zerolog-message-field="message"
                                  sets the zerolog message field name, works with
                                  --output=zerolog
          --zerolog-error-field="error"
                                  sets the zerolog error field name, works with
                                  --output=zerolog
    
    Commands:
      help [<command>...]
        Show help.
    
      run*
        Display logs
    
      version
        Display current version`,
  cmd: [
    {
      name: "kail",
      flags: [
        {
          name: "label",
          type: FlagType.Array,
        },
        {
          name: "pod",
          type: FlagType.Array,
        },
        {
          name: "ns",
          type: FlagType.Array,
        },
        {
          name: "ignore-ns",
          type: FlagType.Array,
        },
        {
          name: "svc",
          type: FlagType.Array,
        },
        {
          name: "rc",
          type: FlagType.Array,
        },
        {
          name: "rs",
          type: FlagType.Array,
        },
        {
          name: "ds",
          type: FlagType.Array,
        },
        {
          name: "deploy",
          type: FlagType.Array,
        },
        {
          name: "sts",
          type: FlagType.Array,
        },
        {
          name: "job",
          type: FlagType.Array,
        },
        {
          name: "ing",
          type: FlagType.Array,
        },
        {
          name: "context",
          type: FlagType.String,
        },
        {
          name: "current-ns",
          type: FlagType.Boolean,
        },
        {
          name: "containers",
          type: FlagType.Array,
        },
        {
          name: "dry-run",
          type: FlagType.Boolean,
        },
        {
          name: "log-file",
          type: FlagType.String,
        },
        {
          name: "log-level",
          type: FlagType.String,
        },
        {
          name: "since",
          type: FlagType.String,
        },
        {
          name: "output",
          type: FlagType.Enum,
          options: ["default", "raw", "json", "json-pretty", "zerolog"],
        },
        {
          name: "zerolog-timestamp-field",
          type: FlagType.String,
        },
        {
          name: "zerolog-level-field",
          type: FlagType.String,
        },
        {
          name: "zerolog-message-field",
          type: FlagType.String,
        },
        {
          name: "zerolog-error-field",
          type: FlagType.String,
        },
      ],
    },
  ],
};

const genSchemaComponents = (raw: CLIJson) => {
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
                  style: "overflow:auto",
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
        type: "custom/v1/TextDisplay",
        properties: {
          text: `${raw.help}`,
          format: '{{"code"}}',
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
              styles: [
                {
                  styleSlot: "content",
                  style: "height:300px;",
                  cssProperties: {},
                },
              ],
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
    flag: CLIJson["cmd"][0]["flags"][0]
  ) => {
    const fieldId = getFlagFieldId(cmdName, flag.name);
    const inputId = getFlagInputId(cmdName, flag.name);
    const validationId = getFlagValidationId(cmdName, flag.name);
    const options = flag.options;

    switch (flag.type) {
      case FlagType.Number:
        return {
          id: inputId,
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
                  id: fieldId,
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
      case FlagType.Array:
        return {
          id: inputId,
          type: "custom/v1/arrayInput",
          properties: {
            value: [""],
            type: "string",
            placeholder: "please input",
            disabled: `{{${cmdName}FormState.data.isExecuting}}`,
          },
          traits: [
            {
              type: "core/v2/slot",
              properties: {
                container: {
                  id: fieldId,
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
      case FlagType.Boolean:
        return {
          id: inputId,
          type: "arco/v1/switch",
          properties: {
            defaultChecked: false,
            disabled: `{{${cmdName}FormState.data.isExecuting}}`,
            loading: false,
            type: "circle",
            size: "default",
            updateWhenDefaultValueChanges: false,
          },
          traits: [
            {
              type: "core/v2/slot",
              properties: {
                container: {
                  id: fieldId,
                  slot: "content",
                },
                ifCondition: true,
              },
            },
          ],
        };
      case FlagType.Enum:
        return {
          id: inputId,
          type: "arco/v1/select",
          properties: {
            allowClear: false,
            multiple: false,
            allowCreate: false,
            bordered: true,
            defaultValue: undefined,
            disabled: false,
            labelInValue: false,
            loading: false,
            showSearch: false,
            unmountOnExit: true,
            showTitle: false,
            options: options?.map((op) => ({ label: op, value: op })) || [],
            placeholder: "Select a option",
            size: "default",
            error: false,
            updateWhenDefaultValueChanges: false,
            autoFixPosition: false,
            autoAlignPopupMinWidth: false,
            autoAlignPopupWidth: true,
            autoFitPosition: false,
            position: "bottom",
            mountToBody: true,
          },
          traits: [
            {
              type: "core/v2/slot",
              properties: {
                container: {
                  id: "kailFormOutputField",
                  slot: "content",
                },
                ifCondition: true,
              },
            },
          ],
        };
      default:
        return {
          id: inputId,
          type: "arco/v1/input",
          properties: {
            allowClear: false,
            disabled: `{{${cmdName}FormState.data.isExecuting}}`,
            readOnly: false,
            defaultValue: "",
            updateWhenDefaultValueChanges: false,
            placeholder: "please input",
            error: `{{${cmdName}Form.validatedResult.${validationId}.isInvalid}}`,
            size: "default",
          },
          traits: [
            {
              type: "core/v2/slot",
              properties: {
                container: {
                  id: fieldId,
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

  const genCmdFormFields = (
    cmdName: string,
    flags: CLIJson["cmd"][0]["flags"]
  ) => {
    let formFieldComponents = [] as unknown[];
    flags.forEach((flag) => {
      const argFieldId = getFlagFieldId(cmdName, flag.name);
      const components = [
        {
          id: argFieldId,
          type: "arco/v1/formControl",
          properties: {
            label: {
              format: "plain",
              raw: `${flag.name}`,
            },
            layout: "horizontal",
            required: flag.required || false,
            hidden: false,
            extra: "",
            errorMsg: "",
            labelAlign: "left",
            colon: false,
            help: "",
            labelCol: {
              span: 6,
              offset: 0,
            },
            wrapperCol: {
              span: 18,
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
        genCmdFormFiledInput(cmdName, flag),
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
                validators: item.flags.map((sub) => {
                  const inputId = getFlagInputId(item.name, sub.name);
                  const validationId = getFlagValidationId(item.name, sub.name);
                  const validation = {
                    name: validationId,
                    value: `{{${inputId}.value}}`,
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
        ...genCmdFormFields(item.name, item.flags),
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
          id: `${item.name}TabResultContent`,
          type: "custom/v1/TextDisplay",
          properties: {
            text: "No Result",
            format: "plain",
          },
          traits: [
            {
              type: "core/v2/slot",
              properties: {
                container: {
                  id: `${item.name}TabResult`,
                  slot: "content",
                },
                ifCondition: true,
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

export const genApp = () => {
  const components = genSchemaComponents(kailRaw);
  const app = {
    kind: "Application",
    version: "CLI2UI/v1",
    metadata: {
      name: "PingCLI",
    },
    spec: {
      components,
    },
  } as Application;
  return app;
};
