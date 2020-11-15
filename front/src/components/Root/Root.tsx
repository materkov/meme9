import React from "react";
import {resolveRoute} from "../../RouteResolver";
import {fetchData} from "../../DataFetcher";
import {fetchJs} from "../../JsFetcher";
import {Error} from "../Error/Error";
import {NavigateContext} from "../../context";

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
        fetchJs(window.InitJsBundles).then(() => {
            this.setState({
                rootData: window.InitApiResponse,
                rootComponent: window.InitRootComponent,
            });
        })
    }

    navigate = (url: string) => {
        resolveRoute(url).then((resolvedRoute) => {
            Promise.all([
                fetchJs(resolvedRoute.js || []),
                fetchData<any>(resolvedRoute.apiMethod, resolvedRoute.apiRequest)
            ]).then(([_, renderer]) => {
                this.setState({
                    rootData: renderer,
                    rootComponent: resolvedRoute.rootComponent,
                })
                window.history.replaceState({}, 'meme', url);
            });
        })
    }

    renderRoot() {
        if (!this.state.rootData || !this.state.rootComponent) {
            return '';
        }

        if (this.state.rootData.error) {
            return <Error data={this.state.rootData.error}/>
        }

        const Component = window.modules[this.state.rootComponent];

        //@ts-ignore
        return <Component data={this.state.rootData.data}/>
    }

    render() {
        return (
            <NavigateContext.Provider value={this.navigate}>
                {this.renderRoot()}
            </NavigateContext.Provider>
        );
    }
}
