import React from "react";
import {resolveRoute} from "../../Router";
import {GlobalContext} from "../../Context";
import {Header} from "../Header/Header";

interface State {
    componentName?: string;
    data?: any;
}

export class Router extends React.Component<{}, State> {
    state: State = {};

    componentDidMount() {
        this.navigate(window.location.pathname);
    }

    navigate = (route: string) => {
        window.history.pushState(null, "meme", route);

        resolveRoute(route).then(([component, data]) => {
            this.setState({
                componentName: component,
                data: data.renderer,
            })
        })
    }

    render() {
        if (!this.state.componentName) {
            return null;
        }

        const Component = window.modules[this.state.componentName];

        return (
            <GlobalContext.Provider value={this.navigate}>
                <Header />

                {this.state.data && <Component data={this.state.data}/>}
            </GlobalContext.Provider>
        )
    }
}