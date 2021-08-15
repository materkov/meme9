import React from "react";
import {UniversalRenderer} from "../types"
import {UniversalComponent} from "./UniversalComponent";
import {UrlContext} from "../context";
import {Link} from "./Link";

export const App = () => {
    const [url, setUrl] = React.useState<string>(window.location.hash.substr(1));
    const [data, setData] = React.useState<UniversalRenderer | undefined>(undefined);

    React.useEffect(() => {
        history.pushState(null,'','#' + url);

        fetch("http://localhost:8000" + url).then(r => r.json()).then(r => {
            setData(r);
        });
    }, [url]);

    if (!data) {
        return <div>Loading...</div>
    }

    return (
        <UrlContext.Provider value={{url: url, navigate: setUrl}}>
            MENU:&nbsp;&nbsp;
            <Link href={"/feed"}>Feed</Link> | <Link href={"/new_post"}>New post</Link> | <Link href={"/vk"}>Login</Link><br/><br/>

            <UniversalComponent data={data}/>
        </UrlContext.Provider>
    );
};
