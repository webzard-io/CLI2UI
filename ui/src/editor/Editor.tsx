import { initSunmaoUIEditor } from "@sunmao-ui/editor";
import {
  getLibs,
  BaseProps,
  saveApp,
  saveModules,
  patchApp,
  patchModules,
  fetchApp,
  fetchModules,
} from "../shared";
import { APPLICATION_NAME } from "../constants";
import "@sunmao-ui/arco-lib/dist/index.css";
import "@sunmao-ui/editor/dist/index.css";
import { type Application, type Module } from "@sunmao-ui/core";
import { useState, useEffect } from "react";
import { genApp } from "../application/utils";
import { formatCommand } from "../sunmao/format-command";

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
      if (application) {
        setInitialized(true);
        return;
      }

      // only fetch application schema when not provide from HTML
      const [appPatch, modulesPatch] = await Promise.all([
        fetchApp(APPLICATION_NAME),
        fetchModules(),
      ]);

      setApp(patchApp(genApp(), appPatch));
      setModules(patchModules(modules, modulesPatch));
      setInitialized(true);
    })();
  }, [application]);

  if (!initialized) return null;

  const { Editor } = initSunmaoUIEditor({
    defaultApplication: _app,
    defaultModules: _modules,
    runtimeProps: {
      libs: getLibs({ ws, handlers, utilMethods }),
      dependencies: {
        // formatCommand,
      },
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
