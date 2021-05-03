import React from "react";
import {resolveRoute} from "../../Router";
import {GlobalContext} from "../../Context";
import {Header} from "../Header/Header";
import * as schema from "../../api/renderer";
import {UniversalRenderer} from "../UniversalRenderer/UniversalRenderer";

interface State {
    data?: schema.UniversalRenderer;
    error?: boolean;
}

export class Router extends React.Component<{}, State> {
    state: State = {};

    componentDidMount() {
        this.navigate(window.location.pathname);
    }

    navigate = (route: string) => {
        window.history.pushState(null, "meme", route);

        resolveRoute(route)
            .then(data => this.setState({data: data, error: undefined}))
            .catch(() => this.setState({data: undefined, error: true}))
    }

    render() {
        return (
            <GlobalContext.Provider value={this.navigate}>
                <Header/>
                {this.state.data && <UniversalRenderer data={this.state.data}/>}
                {this.state.error && <div>Ошибка!</div>}
            </GlobalContext.Provider>
        )
    }
}