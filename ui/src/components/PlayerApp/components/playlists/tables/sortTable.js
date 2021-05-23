import "./style.less"
import React, {useState, useEffect} from "react"
import { Table } from 'antd';
import { sortableContainer, sortableElement } from 'react-sortable-hoc';
import arrayMove from 'array-move';
import {SortableCols} from "./colomns"

const SortableItem = sortableElement(props => <tr {...props} />);
const SortableContainer = sortableContainer(props => <tbody {...props} />);

export const SortableTable = ({data, setData}) => {
  const [datastore, setDatastore] = useState(data)
  useEffect(() => {setDatastore(data)}, [data]);
  const onSortEnd = ({ oldIndex, newIndex }) => {
    if (oldIndex !== newIndex) {
      const newData = arrayMove([].concat(datastore), oldIndex, newIndex).filter(el => !!el);
      console.log('Sorted items: ', newData);
      setData(newData)
    }
  }
  const DraggableContainer = props => (
    <SortableContainer
      useDragHandle
      disableAutoscroll
      helperClass="row-dragging"
      onSortEnd={onSortEnd}
      {...props}
    />
  )
  const DraggableBodyRow = ({ className, style, ...restProps }) => {
    // function findIndex base on Table rowKey props and should always be a right array index
    const index = datastore.findIndex(x => x.id === restProps['data-row-key']);
    return (<SortableItem className="editable-row" index={index} {...restProps} />)
  }
  return (
        <Table
            style={{height:"500px", "overflowY":"auto"}}
            pagination={false}
            dataSource={datastore}
            columns={SortableCols()}
            rowKey="id"
            components={{
                body: {
                wrapper: DraggableContainer,
                row: DraggableBodyRow,
                },
            }}
        />
    );
}
