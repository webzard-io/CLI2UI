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
} from "../shared";
import { APPLICATION_NAME } from "../constants";
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

  return <div>{initialized ? <SunmaoApp options={_app} /> : null}</div>;
}

export default App;
