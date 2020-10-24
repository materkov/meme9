import React from "react";
import * as schema from "../../schema/login";
import {resolveRoute} from "../../RouteResolver";
import {fetchData} from "../../DataFetcher";
import {loadJs} from "../../JsManager";

interface State {
    rootData: schema.AnyRenderer | undefined;
}

export class Root extends React.Component<{}, State> {
    state: State = {
        rootData: undefined,
    };

    componentDidMount() {
        if (window.InitData) {
            this.setState({rootData: window.InitData});
        }
    }

    onClick = (e: React.MouseEvent) => {
        //@ts-ignore
        if (!e.target || e.target.nodeName !== "A") {
            return;
        }

        e.preventDefault();

        //@ts-ignore
        const url = e.target.href;

        resolveRoute(url).then((resolvedRoute) => {
            if (!resolvedRoute.request) {
                return
            }

            Promise.all([
                loadJs(resolvedRoute.js),
                fetchData(resolvedRoute.request)
            ]).then(([_, renderer]) => {
                this.setState({rootData: renderer})
                window.history.replaceState({}, 'meme', url);
            })
        })
    }

    renderRoot() {
        if (!this.state.rootData) {
            return '';
        }

        let key = ""
        for (key in this.state.rootData) {
            break;
        }
        const componentName = key[0].toUpperCase() + key.substring(1);
        const Component = window.modules[componentName];

        //@ts-ignore
        return <Component {...this.state.rootData[key]}/>
    }

    render() {
        return <div onClick={this.onClick} id="content">{this.renderRoot()}</div>;
    }
}
