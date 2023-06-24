export class Article {
    type = Article.name
    id = ""
    title = ""
    user: User | null = null
    createdAt = ""
    paragraphs: Paragraph[] = []

    static fromJSON(raw: any): Article {
        const o = new Article();
        o.id = String(raw.id || "")
        o.title = String(raw.title || "")
        o.createdAt = String(raw.createdAt || "")
        o.paragraphs = (raw.paragraphs || []).map(Paragraph.fromJSON)
        o.user = raw.user ? User.fromJSON(raw.user) : null
        return o;
    }
}

export class Paragraph {
    text: ParagraphText | undefined
    image: ParagraphImage | undefined
    list: ParagraphList | undefined

    static fromJSON(raw: any): Paragraph {
        const o = new Paragraph();
        if (raw.text) {
            o.text = ParagraphText.fromJSON(raw.text);
        }
        if (raw.image) {
            o.image = ParagraphImage.fromJSON(raw.image);
        }
        return o
    }
}

export class ParagraphText {
    type = ParagraphText.name
    id = ""
    text = ""

    static fromJSON(raw: any): ParagraphText {
        const o = new ParagraphText();
        o.id = String(raw.id || "")
        o.text = String(raw.text || "")
        return o
    }
}

export class ParagraphImage {
    type = ParagraphImage.name
    id = ""
    url = ""

    static fromJSON(raw: any): ParagraphImage {
        const o = new ParagraphImage();
        o.id = String(raw.id || "")
        o.url = String(raw.url || "")
        return o
    }
}

export class ParagraphList {
    id = ""
    type: ListType = ListType.UNKNOWN
    items: string[] = []
}

export enum ListType {
    UNKNOWN = "",
    ORDERED = "ORDERED",
    UNORDERED = "UNORDERED",
}

export class User {
    id = ""
    name = ""

    static fromJSON(raw: any): User {
        const o = new User();
        o.id = String(raw.id || "")
        o.name = String(raw.name || "")
        return o
    }
}

export class ArticlesList {
    id: string = ""
}

export class ArticlesSave {
    id: string = ""
    title: string = ""
    paragraphs: InputParagraph[] = []
}

export class InputParagraph {
    inputParagraphText: InputParagraphText | undefined = undefined
    inputParagraphImage: InputParagraphImage | undefined = undefined
}

export class InputParagraphText {
    text: string = ""
}

export class InputParagraphImage {
    url: string = ""
}
