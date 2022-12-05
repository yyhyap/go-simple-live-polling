import React from 'react';
import { useEffect } from 'react';
import { useState } from 'react';
import { partiesData } from './Data';
import PartyItem from './PartyItem';

const Parties = () => {

    const baseURL = 'ws://localhost:8000/ws/live-vote';
    const [parties, setParties] = useState([]);
    const [websocket, setWebsocket] = useState(undefined);
    const [data, setData] = useState(undefined);

    useEffect(() => {
        if(websocket === undefined) {
            let ws = new WebSocket(baseURL);

            ws.onopen = (e) => {
                console.log('Websocket connection established!', {e});
            }

            ws.onclose = (e) => {
                console.log('Websocket connection closed!', {e});
                setWebsocket(undefined);
            }

            ws.onmessage = (msg) => {
                console.log('Websocket message: ', {msg});
                let rawData = JSON.parse(msg.data);
                const newData = new Map();
                rawData.data.forEach((item) => {
                    newData.set(item._id, item.totalCount);
                })
                setData(newData);
            }

            ws.onerror = (error) => {
                console.log('Websocket error: ', {error});
                setWebsocket(undefined);
            }

            setWebsocket(ws);
        }
    }, [websocket])

    useEffect(() => {
        if(data !== undefined) {
            const newPartiesData = [...parties];
            newPartiesData.forEach((item) => {
                if(data.has(item.id)) {
                    item.totalCount = data.get(item.id);
                }
            })
            setParties(newPartiesData);
        } else {
            setParties(partiesData);
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [data])

    return (
        <div>
            <div className='party__container container grid'>
                {parties.map((item) => {
                    return (
                        <PartyItem item={item} key={item.id} />
                    )
                })}
            </div>
        </div>
    )
}

export default Parties