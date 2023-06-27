import React, {useEffect} from "react";
import {Article} from "./Paragraphs/Article";
import {useArticlePage} from "../store/articlePage";
import * as styles from "./ArticlePage.module.css";

export function ArticlePage() {
    let articleId = window.document.location.pathname.substring(9);
    const articleState = useArticlePage()

    useEffect(() => {
        articleState.fetch(articleId);
    }, []);

    return (
        <div>
            {!articleState.error ?
                <Article articleId={articleId}/> :
                <div className={styles.error}>error</div>
            }
        </div>
    )
}
