import * as schema from "../../schema/login";
import React from "react";

interface ErrorProps {
    data: schema.ErrorRenderer;
}

export function Error(props: ErrorProps) {
    return <div style={{fontSize: '20px'}}>
        Произошла ошибочка: <b>{props.data.displayText}</b>
    </div>
}
