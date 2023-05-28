import React, {useEffect} from "react";
import {useProfile} from "../../store/profile";
import {ArticleSnippet} from "../ArticleSnippet/ArticleSnippet";

export function Profile() {
    const userId = document.location.pathname.substring(7);
    const profileState = useProfile();

    useEffect(() => {
        profileState.fetch(userId);
    }, []);

    if (!profileState.user.id) {
        return <div>Loading....</div>
    }

    return <div>
        <h1>{profileState.user.name}</h1>
        <hr/>

        {profileState.articles.map(article => (
            <ArticleSnippet article={article} key={article.id}/>
        ))}
    </div>
}


