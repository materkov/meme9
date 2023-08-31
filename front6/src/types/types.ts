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

export class PostsAddReq {
    text: string = ""
}

export class Post {
    id: string = ""
    userId: string = ""
    date: string = ""
    text: string = ""
    user?: User = undefined
}

export class PostsListPostedByUser {
    userId: string = ""
}

export class PostsListById {
    id: string = ""
}

export class Void {
}
