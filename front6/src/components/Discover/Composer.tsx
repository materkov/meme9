import * as styles from "./Composer.module.css";
import React from "react";
import {useQueryClient} from "@tanstack/react-query";
import {uploadFile} from "../../api/uploads";
import {ApiPolls, ApiPosts} from "../../api/client";

export function Composer() {
    const [text, setText] = React.useState('');
    const [saving, setSaving] = React.useState(false);
    const queryClient = useQueryClient();

    const post = () => {
        if (!text) {
            return;
        }

        setSaving(true);
        ApiPosts.Add({text: text, pollId: pollId, photoId: photoId}).then(() => {
            setSaving(false);
            setText('');
            setPhotoId('');
            queryClient.invalidateQueries({queryKey: ['discover']});
        })
    };

    const [question, setQuestion] = React.useState('');
    const [answers, setAnswers] = React.useState<string[]>(['']);

    const savePoll = () => {
        ApiPolls.Add({question: question, answers: answers}).then(poll => {
            setPollId(poll.id);
            setPollActive(false);
            setQuestion('');
            setAnswers([]);
        })
    }

    const addPoll = () => {
        const input = document.createElement('input');
        input.type = 'file';
        input.click();
        // @ts-ignore
        input.onchange = (e: React.ChangeEvent<HTMLInputElement>) => {
            if (!e.target.files) return;

            uploadFile(e.target.files[0]).then(photoId => setPhotoId(photoId));
        };
    }

    const [pollActive, setPollActive] = React.useState(false);
    const [pollId, setPollId] = React.useState('');
    const [photoId, setPhotoId] = React.useState('');

    return <>
        <div className={styles.newPostContainer}>
                <textarea className={styles.newPost} placeholder="What's new today?" value={text}
                          onChange={e => setText(e.target.value)}/>

            {pollId && <>Poll attached</>}

            {!pollId && !pollActive && <a href="#" onClick={(e) => {
                e.preventDefault();
                setPollActive(true);
            }}>Add poll</a>}

            &nbsp;|&nbsp;
            {!photoId &&
                <a href="#" onClick={(e) => {
                    e.preventDefault();
                    addPoll();
                }}>Add photo</a>
            }
            {photoId && <span>Photo attached</span>}

            {pollActive && <>
                <input className={styles.pollAnswer} value={question} onChange={e => setQuestion(e.target.value)}
                       type="text" placeholder={"Your question"}/>

                {answers.map((answer, i) => (
                    <input className={styles.pollAnswer} type="text" placeholder={"Answer option..."} value={answer}
                           onChange={e => {
                               const answersCopy = structuredClone(answers);
                               answersCopy[i] = e.target.value;
                               setAnswers(answersCopy);
                           }}/>
                ))}

                <button onClick={() => setAnswers([...answers, ''])}>Add answer</button>
                <br/>
                <button onClick={savePoll}>Create poll</button>
            </>
            }
            <hr/>

            <button disabled={saving} onClick={post}>Post</button>
            <hr/>
        </div>
    </>
}