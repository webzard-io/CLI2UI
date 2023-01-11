import { initSunmaoUI } from "@sunmao-ui/runtime";
import "@sunmao-ui/arco-lib/dist/index.css";
import {
  getLibs,
  useApiService,
  BaseProps,
  patchApp,
  patchModules,
  fetchApp,
  fetchModules,
  APPLICATION_NAME,
} from "../shared";
import { RuntimeModule } from "@sunmao-ui/core";
import { type Application, type Module } from "@sunmao-ui/core";
import { useState, useEffect } from "react";
import { genApp } from "./utils";

function App(props: BaseProps) {
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

  useEffect(() => {
    (async function () {
      const [appPatch, modulesPatch] = await Promise.all([
        fetchApp(APPLICATION_NAME),
        fetchModules(),
      ]);
      setApp(patchApp(genApp(), appPatch));
      setModules(patchModules(modules, modulesPatch));
    })();
  }, []);

  const {
    App: SunmaoApp,
    apiService,
    registry,
  } = initSunmaoUI({
    libs: getLibs({ ws, handlers, utilMethods }),
  });

  if (_modules) {
    _modules.forEach((moduleSchema) => {
      registry.registerModule(moduleSchema as RuntimeModule);
    });
  }

  useApiService({ ws, apiService });

  return <SunmaoApp options={_app} />;
}

export default App;
