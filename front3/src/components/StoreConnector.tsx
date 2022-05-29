import React from "react";
import {State} from "../state";
import {RootRenderer} from "./types";
import {Root} from "./Root";

export const StateContext = React.createContext<State>(new State());

export const StoreConnector = () => {
    const [rootRenderer, setRootRenderer] = React.useState<RootRenderer | null>(null);
    const [state, setState] = React.useState<State>(new State());

    React.useEffect(() => {
        const state = new State();
        state.subscribe((renderer: RootRenderer) => {
            setRootRenderer({...renderer});
        })
        setState(state);
        setRootRenderer(state.rootRenderer);
    }, []);

    return (
        <StateContext.Provider value={state}>
            {rootRenderer && <Root data={rootRenderer}/>}
        </StateContext.Provider>
    )
}
