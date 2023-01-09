import { initSunmaoUIEditor } from "@sunmao-ui/editor";
import {
  getLibs,
  BaseProps,
  saveApp,
  saveModules,
  // patchApp,
  // patchModules,
  fetchApp,
  fetchModules,
  APPLICATION_NAME,
} from "../shared";
import "@sunmao-ui/arco-lib/dist/index.css";
import "@sunmao-ui/editor/dist/index.css";
import { type Application, type Module } from "@sunmao-ui/core";
import { useState, useMemo, useEffect } from "react";
import { genSchemaComponents } from "../application/utils";

function Editor(props: BaseProps) {
  const {
    application,
    modules,
    handlers,
    ws,
    utilMethods,
    // applicationPatch,
    // modulesPatch,
  } = props;
  const [_app, setApp] = useState<Application>(application);
  const [_modules, setModules] = useState<Module[]>(modules);
  const [initialized, setInitialized] = useState(false);

  useEffect(() => {
    (async function () {
      const [app, modules] = await Promise.all([
        fetchApp(APPLICATION_NAME),
        fetchModules(),
      ]);

      const components = genSchemaComponents();
      app.spec.components = components;
      setApp(app);
      setModules(modules);
      setInitialized(true);
    })();
  }, []);

  if (!initialized) return null;

  const { Editor } = initSunmaoUIEditor({
    // defaultApplication: patchApp(application, applicationPatch),
    // defaultModules: patchModules(modules, modulesPatch)
    defaultApplication: _app,
    defaultModules: _modules,
    runtimeProps: {
      libs: getLibs({ ws, handlers, utilMethods }),
    },
    storageHandler: {
      onSaveApp: function (newApp) {
        saveApp(newApp, _app);
      },
      onSaveModules: function (newModules) {
        saveModules(newModules, _modules || []);
      },
    },
  });

  // TODO: call the useApiService hook when sunmao-ui expose apiService in editor mode

  return <Editor />;
}

export default Editor;
