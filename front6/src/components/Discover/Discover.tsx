import React, {useEffect} from "react";
import {useDiscoverPage} from "../../store/discoverPage";
import {ArticleSnippet} from "../ArticleSnippet/ArticleSnippet";

export function Discover() {
    const discoverState = useDiscoverPage();

    useEffect(() => {
        discoverState.fetch();
    }, []);

    return <div>
        <h1>Discover</h1>
        {discoverState.articles.map(article => <ArticleSnippet article={article} key={article.id}/>)}
    </div>
}
