import { implementUtilMethod } from "@sunmao-ui/runtime";
import React from "react";
import ReactDOM from "react-dom";
import App from "./App";
import { MainOptions } from "../shared";

export function renderApp(options: MainOptions) {
  const {
    wsUrl,
    application,
    modules,
    reloadWhenWsDisconnected,
    handlers,
    utilMethods,
    applicationPatch,
    modulesPatch,
  } = options;
  const ws = new WebSocket(wsUrl);
  ws.onopen = () => {
    console.log("ws connected");
  };
  ws.onclose = () => {
    if (reloadWhenWsDisconnected) {
      setTimeout(() => {
        window.location.reload();
      }, 1500);
    }
  };

  ws.addEventListener('message', e => {
    const data = JSON.parse(e.data)
    
    if (data.handler !== 'Heartbeat') {
      return
    }

    ws.send(JSON.stringify({
      type: 'Action',
      handler: 'Heartbeat',
      // TODO(xinxi.guo): this can be extended later
      params: "Pong",
      store: null
    }))
  })

  ReactDOM.render(
    <App
      application={application}
      modules={modules}
      ws={ws}
      handlers={handlers}
      utilMethods={utilMethods?.map(
        (u) => () => implementUtilMethod(u.options)(u.impl)
      )}
      applicationPatch={applicationPatch}
      modulesPatch={modulesPatch}
    />,
    document.getElementById("root")!
  );
}
