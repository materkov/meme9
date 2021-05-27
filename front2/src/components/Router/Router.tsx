import React from "react";
import {GlobalStoreContext} from "../../Context";
import {Header} from "../Header/Header";
import * as schema from "../../api/renderer";
import {UniversalRenderer} from "../UniversalRenderer/UniversalRenderer";
import {Store} from "../../Store";
import {HeaderRenderer} from "../../api/api2";

interface State {
    data?: schema.UniversalRenderer;
    headerData: HeaderRenderer;
    error?: boolean;
}

export class Router extends React.Component<{}, State> {
    state: State;
    globalStore: Store;

    constructor(props: any) {
        super(props);

        this.globalStore = new Store((d: schema.UniversalRenderer, d2: HeaderRenderer) => {
            this.setState({data: d, headerData: d2});
        });
        this.state = {
            headerData: this.globalStore.headerData,
            data: this.globalStore.data,
        }
    }

    render() {
        return (
            <GlobalStoreContext.Provider value={this.globalStore}>
                <Header data={this.state.headerData}/>
                {this.state.data && <UniversalRenderer data={this.state.data}/>}
                {this.state.error && <div>Ошибка!</div>}
            </GlobalStoreContext.Provider>
        )
    }
}
