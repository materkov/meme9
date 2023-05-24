import React, {useEffect} from "react";
import {Article} from "./Paragraphs/Article";
import {ArticlePage, getArticle, useArticlePage} from "../store/store";
import * as styles from "./App.module.css";

export function App() {
    let articleId = window.document.location.pathname;
    if (articleId.indexOf("/article/") == 0) {
        articleId = articleId.substring(9);
    }

    const [error, setError] = React.useState("");
    const article = useArticlePage((state: ArticlePage) => state.article)
    const setArticle = useArticlePage((state: ArticlePage) => state.setArticle);

    useEffect(() => {
        getArticle(articleId)
            .then(setArticle)
            .catch((error) => {
                if (error.code == 404) {
                    setError("Статья не существует или была удалена")
                } else {
                    setError("Не удалось загрузить статью")
                }
            })
    }, [])

    return (
        <div>
            {article ?
                <Article/> :
                <div className={styles.error}>{error}</div>
            }
        </div>
    )
}
