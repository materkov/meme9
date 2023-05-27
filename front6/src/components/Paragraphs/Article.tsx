import {ArticlesSave, InputParagraph, InputParagraphImage, InputParagraphText} from "../../types/types";
import {Image} from "./Image";
import {Text} from "./Text";
import React from "react";
import * as styles from "./Article.module.css";
import {ArticlePage, useArticlePage} from "../../store/articlePage";
import {formatDate} from "../../utils/date";
import {Link} from "../Link";
import {api} from "../../store/api";

export function Article() {
    const article = useArticlePage((state: ArticlePage) => state.article);

    const items = article.paragraphs.map(p => {
        if (p.image) {
            return <Image key={p.image.id} {...p.image}/>
        } else if (p.text) {
            return <Text key={p.text.id} {...p.text}/>;
        } else {
            return null;
        }
    })

    const onSave = () => {
        const req: ArticlesSave = {
            id: article.id,
            title: article.title,
            paragraphs: article.paragraphs.map(p => {
                const outer = new InputParagraph();
                if (p.text) {
                    outer.inputParagraphText = new InputParagraphText();
                    outer.inputParagraphText.text = p.text.text;
                } else if (p.image) {
                    outer.inputParagraphImage = new InputParagraphImage();
                    outer.inputParagraphImage.url = p.image.url;
                }

                return outer;
            })
        };

        api<void>("articles.save", req)
    }

    return <div className={styles.article}>
        <h1 className={styles.title}>{article.title}</h1>

        <div className={styles.author}>
            <Link href={"/users/" + article.user?.id}>
                {article.user?.name}
            </Link>
            <div className={styles.time}>{formatDate(article.createdAt)}</div>
        </div>

        {items}
        <button onClick={onSave}>Save</button>
    </div>;
}
