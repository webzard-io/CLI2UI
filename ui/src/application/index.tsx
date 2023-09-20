import { implementUtilMethod } from "@sunmao-ui/runtime";
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

  ws.addEventListener('open', () => {
    ws.send(JSON.stringify({
      type: 'Action',
      handler: 'ConnectionEstablished',
      params: {
        clientId: localStorage.getItem('clientId'),
        serverSignature: localStorage.getItem('serverSignature')
      },
      store: null
    }))
  });

  ws.addEventListener('close', () => {
    if (reloadWhenWsDisconnected) {
      setTimeout(() => {
        window.location.reload();
      }, 1500);
    }
  });

  ws.addEventListener('message', e => {
    const data = JSON.parse(e.data)
    
    if (data.handler === 'Heartbeat') {
      if (data.params.hasOwnProperty('serverSignature')) {
        localStorage.setItem('clientId', data.params.clientId)
        localStorage.setItem('serverSignature', data.params.serverSignature)
        // FIXME(xinxi.guo): this is a workaround, right after clientId is set, websocket does not work as expected
        window.location.reload();
      }

      ws.send(JSON.stringify({
        type: 'Action',
        handler: 'Heartbeat',
        params: "Pong",
        store: null
      }))
    }
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
