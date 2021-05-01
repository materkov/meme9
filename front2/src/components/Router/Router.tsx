import React from "react";
import {resolveRoute} from "../../Router";
import {GlobalContext} from "../../Context";
import {Header} from "../Header/Header";
import * as schema from "../../api/renderer";
import {UniversalRenderer} from "../UniversalRenderer/UniversalRenderer";

interface State {
    data?: schema.UniversalRenderer;
}

export class Router extends React.Component<{}, State> {
    state: State = {};

    componentDidMount() {
        this.navigate(window.location.pathname);
    }

    navigate = (route: string) => {
        window.history.pushState(null, "meme", route);

        resolveRoute(route).then(data => {
            this.setState({data: data})
        })
    }

    render() {
        if (!this.state.data) {
            return null;
        }

        return (
            <GlobalContext.Provider value={this.navigate}>
                <Header/>
                <UniversalRenderer data={this.state.data}/>
            </GlobalContext.Provider>
        )
    }
}