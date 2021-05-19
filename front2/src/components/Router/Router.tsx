import React from "react";
import {GlobalStoreContext} from "../../Context";
import {Header} from "../Header/Header";
import * as schema from "../../api/renderer";
import {UniversalRenderer} from "../UniversalRenderer/UniversalRenderer";
import {Store} from "../../Store";

interface State {
    data?: schema.UniversalRenderer;
    error?: boolean;
}

export class Router extends React.Component<{}, State> {
    state: State = {};
    globalStore: Store;

    constructor(props: any) {
        super(props);

        this.globalStore = new Store((d: schema.UniversalRenderer) => {
            this.setState({data: d});
        });
    }

    componentDidMount() {
        this.globalStore.navigate(window.location.pathname);
    }

    render() {
        return (
            <GlobalStoreContext.Provider value={this.globalStore}>
                <Header/>
                {this.state.data && <UniversalRenderer data={this.state.data}/>}
                {this.state.error && <div>Ошибка!</div>}
            </GlobalStoreContext.Provider>
        )
    }
}
