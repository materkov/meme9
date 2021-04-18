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
                <Header data={{
                    mainUrl: "/",
                    userName: "Макс",
                    userAvatar: "https://sun2.43222.userapi.com/s/v1/ig2/WVgdYwZ6Cd8mMcMunD_Po2YDv0_2BHRGf3ofZ1NHyGcbd9nKDQJ029FOYIgwo614Rqv3RT1hO7z5t01DUSRaJosq.jpg?size=100x0&quality=96&crop=122,105,561,561&ava=1",
                }}/>

                {this.state.data && <Component data={this.state.data}/>}
            </GlobalContext.Provider>
        )
    }
}