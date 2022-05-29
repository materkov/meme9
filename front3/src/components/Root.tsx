import {RootRenderer} from "./types";
import React from "react";
import {Header} from "./Header";
import {StateContext} from "./StoreConnector";
import {Composer} from "./Composer";

export const Root = (props: { data: RootRenderer }) => {
    const state = React.useContext(StateContext);
    const onNameClick = () => {
        state.do(props.data.command);
    }

    return <>
        <Header data={props.data.header}/>
        <h1 onClick={onNameClick}>{props.data.label}</h1>
        <Composer data={props.data.composer}/>
    </>
}
