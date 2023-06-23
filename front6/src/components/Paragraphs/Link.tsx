import * as types from "../../types/types";
import {ListType} from "../../types/types";
import React from "react";

export function List(paragraph: types.ParagraphList) {
    const items = paragraph.items.map(item => <li>{item}</li>);

    if (paragraph.type == ListType.ORDERED) {
        return <ul>{items}</ul>
    } else if (paragraph.type == ListType.UNORDERED) {
        return <ol>{items}</ol>;
    } else {
        return null;
    }
}
