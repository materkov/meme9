import React from "react";
import * as schema from "../../schema/login";

interface State {
    rootData: schema.AnyRenderer | undefined;
    loadedJs: { [fileName: string]: boolean };
    routeCache: { [url: string]: schema.AnyRequest };
    dataCache: { [url: string]: schema.AnyRenderer };
}

export class Root extends React.Component<{}, State> {
    state: State = {
        loadedJs: {},
        rootData: undefined,
        routeCache: {},
        dataCache: {},
    };

    componentDidMount() {
        if (window.InitData) {
            this.setState({
                rootData: window.InitData,
                dataCache: {...this.state.dataCache, [window.location.pathname]: window.InitData},
                routeCache: {...this.state.routeCache, [window.location.pathname]: window.InitRequest},
            });
        }

        if (window.InitJsBundles) {
            let loadedJS: { [fileName: string]: boolean } = {};
            for (let js of window.InitJsBundles) {
                loadedJS[js] = true
            }

            this.setState({loadedJs: loadedJS});
        }
    }

    loadAllJs = (js: string[]): Promise<null> => {
        return new Promise<null>((resolve, reject) => {
            let neededJs = [];
            for (let file of js) {
                if (!this.state.loadedJs[file]) {
                    neededJs.push(file);
                }
            }

            if (neededJs.length === 0) {
                resolve();
            }

            let loaded = 0;
            let needed = neededJs.length;

            for (let file of neededJs) {
                const fileName = file;
                const script = document.createElement('script');
                script.src = fileName;
                script.onload = () => {
                    this.setState({loadedJs: {...this.state.loadedJs, [fileName]: true}});

                    loaded++;
                    if (loaded === needed) {
                        resolve();
                    }
                };

                document.body.appendChild(script);
            }
        });
    }

    onClick = (e: React.MouseEvent) => {
        //@ts-ignore
        if (!e.target || e.target.nodeName !== "A") {
            return;
        }

        e.preventDefault();
        //@ts-ignore
        const url = new URL(e.target.href);

        const cachedData = this.state.dataCache[url.pathname];
        if (cachedData) {
            window.history.replaceState({}, 'meme', url.toString());
            this.setState({rootData: cachedData});
            return;
        }

        fetch("/resolve-route", {
            method: 'POST', body: JSON.stringify({url: url.pathname})
        }).then(r => r.json()).then((r: schema.ResolveRouteResponse) => {
            this.loadAllJs(r.js).then(() => {
                if (r.request) {
                    this.setState({
                        routeCache: {...this.state.routeCache, [url.pathname]: r.request},
                    });
                }

                fetch("/api", {
                    method: 'POST',
                    body: JSON.stringify(r.request)
                }).then((r => r.json())).then(r => {
                    this.setState({
                        dataCache: {...this.state.dataCache, [url.pathname]: r},
                        rootData: r,
                    })
                    window.history.replaceState({}, 'meme', url.toString());
                });
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
