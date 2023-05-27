import React, {useEffect} from "react";
import {useProfile} from "../../store/profile";
import * as styles from "./Profile.module.css";
import * as types from "../../types/types";
import {Link} from "../Link";

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
            <ArticleSnippet articleId={article.id} key={article.id}/>
        ))}
    </div>
}

function ArticleSnippet(props: { articleId: string }) {
    const article = useProfile(state => state.articles.find(article => article.id === props.articleId));
    if (!article) {
        return null;
    }

    return <div>
        <Link className={styles.articleTitle} href={"/article/" + article.id}>{article.title}</Link>
        <div className={styles.articleSnippet}>{getTextSnippet(article)}</div>
    </div>
}

function getTextSnippet(article: types.Article): string {
    for (const p of article.paragraphs) {
        if (p.text) {
            return p.text.text;
        }
    }

    return "";
}
