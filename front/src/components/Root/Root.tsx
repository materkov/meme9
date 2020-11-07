import React from "react";
import {resolveRoute} from "../../RouteResolver";
import {fetchData} from "../../DataFetcher";
import {fetchJs} from "../../JsFetcher";

interface State {
    rootData: any;
    rootComponent: string;
}

export class Root extends React.Component<{}, State> {
    state: State = {
        rootData: undefined,
        rootComponent: '',
    };

    componentDidMount() {
        if (window.InitApiResponse) {
            this.setState({
                rootData: window.InitApiResponse,
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
            Promise.all([
                fetchJs(resolvedRoute.js || []),
                fetchData(resolvedRoute.apiMethod, resolvedRoute.apiRequest)
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
