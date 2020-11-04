import React from "react";
import * as schema from "../../schema/login";
import {resolveRoute} from "../../RouteResolver";
import {fetchData} from "../../DataFetcher";
import {fetchJs} from "../../JsFetcher";

interface State {
    rootData: schema.AnyRenderer | undefined;
    rootComponent: string;
}

export class Root extends React.Component<{}, State> {
    state: State = {
        rootData: undefined,
        rootComponent: '',
    };

    componentDidMount() {
        if (window.InitData) {
            this.setState({
                rootData: window.InitData,
                rootComponent: window.InitRootComponent,
            });
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
                fetchJs(resolvedRoute.js || []),
                fetchData(resolvedRoute.request)
            ]).then(([_, renderer]) => {
                this.setState({
                    rootData: renderer,
                    rootComponent: resolvedRoute.rootComponent,
                })
                window.history.replaceState({}, 'meme', url);
            })
        })
    }

    renderRoot() {
        if (!this.state.rootData || !this.state.rootComponent) {
            return '';
        }

        const Component = window.modules[this.state.rootComponent];

        //@ts-ignore
        return <Component data={this.state.rootData}/>
    }

    render() {
        return <div onClick={this.onClick} id="content">{this.renderRoot()}</div>;
    }
}
