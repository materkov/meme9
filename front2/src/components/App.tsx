import React from "react";
import {UniversalRenderer} from "../types"
import {UniversalComponent} from "./UniversalComponent";
import {UrlContext} from "../context";
import {Link} from "./Link";

export const App = () => {
    const [url, setUrl] = React.useState<string>(window.location.pathname);
    //@ts-ignore
    const [data, setData] = React.useState<UniversalRenderer | undefined>(window.__initialData);

    function navigate(url: string) {
        history.pushState(null, '', url);

        fetch("http://localhost:8000" + url, {
            headers: {
                "x-ajax": "1",
            }
        }).then(r => r.json()).then(r => {
            setData(r);
        });
    }

    if (!data) {
        return <div>Loading...</div>
    }

    return (
        <UrlContext.Provider value={{url: url, navigate: navigate}}>
            MENU:&nbsp;&nbsp;
            <Link href={"/feed"}>Feed</Link> | <Link href={"/new_post"}>New post</Link> | <Link
            href={"/vk"}>Login</Link><br/><br/>

            <UniversalComponent data={data}/>
        </UrlContext.Provider>
    );
};
