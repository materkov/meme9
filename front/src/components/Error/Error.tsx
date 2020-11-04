import * as schema from "../../schema/login";
import React from "react";

interface ErrorProps {
    data: schema.ErrorRenderer;
}

function Error(props: ErrorProps) {
    return <div style={{fontSize: '20px'}}>
        {props.data.displayText}
    </div>
}
