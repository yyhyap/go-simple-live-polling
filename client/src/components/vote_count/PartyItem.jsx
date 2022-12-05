import React from 'react'

const PartyItem = ({ item }) => {

    let partyName = item.id;
    partyName = partyName.replace(/([A-Z])/g, ' $1').trim();

    return (
        <div className={`party__card ${item.id}`} key={item.id}>
            <div className='party__id'>
                {partyName}
            </div>
            <div className='party__count'>
                {item.totalCount}
            </div>            
        </div>
    )
}

export default PartyItem