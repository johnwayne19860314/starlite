import  { useState, FC,useEffect } from 'react';
import {

  Button,
  Popconfirm,
  Pagination,
  message,
} from 'antd';
import ProTable, { ProColumns } from '@ant-design/pro-table';
//import {  Dispatch } from 'umi';
//import {roles} from './components/UserModal';
import UserModal , {roles,FormValues} from './components/UserModal';
import {initUserInfo} from '@/models/data';



import {
  userList, addUser, updateUser, delUser,
  
} from '@/services/user';

// interface UserPageProps {
//   users: UserState;
//   dispatch: Dispatch;
//   userListLoading: boolean;
// }
// const initUserInfo = {
//   name: "string",
//   avatar: "string",
//   key: "string",
//   id: 0,
//   email: "string",
//   role:"string"
// }
//const pageSize = 1
const UserListPage: FC = () => {
  const [modalVisible, setModalVisible] = useState(false);
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [record, setRecord] = useState<API.UserInfo >(initUserInfo);
  const [dataSource, setDataSource] = useState<API.UserInfo[]>([]);
  //const [page, setPage] = useState(1)
  const [resetDataSource, setResetDataSource]= useState(false)

  useEffect(() => {
    // Fetch data from your API endpoint
    const fetchData = async () => {
      try {
        const response = await userList({
          page: 0,
          pageSize:100
        }); // Replace with your API endpoint
        setDataSource(response.data.users);
        //setPage(response.data.page)
        console.log(response)
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    fetchData();
  }, [resetDataSource]);

  const columns: ProColumns<API.UserInfo>[] = [
    // {
    //   title: 'ID',
    //   dataIndex: 'id',
    //   key: 'id',
    // },
    {
      title: '名字',
      dataIndex: 'name',
      key: 'name',
      valueType: 'text',
      //render: (text: any) => {text},
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: '角色',
      dataIndex: 'role',
      key: 'role',
      render: (role: any) => <p>{roles[role]}</p>,
    },
    // {
    //   title: 'Create Time',
    //   dataIndex: 'create_time',
    //   valueType: 'dateTime',
    //   key: 'create_time',
    // },
    {
      title: '操作',
      key: 'action',
      valueType: 'option',
      render: (text: any, record: API.UserInfo) => [
        <a
          onClick={() => {
            editHandler(record);
          }}
        >
          更改
        </a>,
        <Popconfirm
          title="确定要删除此用户吗？"
          onConfirm={() => {
            deleteHandler(record.id);
          }}
          okText="确认"
          cancelText="取消"
        >
          <a>删除</a>
        </Popconfirm>,
      ],
    },
  ];

  const editHandler = (record: API.UserInfo) => {
    setModalVisible(true);
    setRecord(record);
  };

  const deleteHandler = (id: number) => {
    const delUserRecord = async () => {
      try {
        const response = await delUser(id); // Replace with your API endpoint
       
        console.log(response)
        resetHandler()
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    delUserRecord();
  };

  const closeHandler = () => {
    setModalVisible(false);
  };

  const onFinish = async (values: FormValues) => {
    //const values = await form.validateFields();
    console.log("=======", values)
    setConfirmLoading(true);
    let id = 0;

    if (record) {
      id = record.id;
    }

    let serviceFun;
    if (id) {
      serviceFun = updateUser;
    } else {
      serviceFun = addUser;
    }
    // name?: string;
    // avatar?: string;
    // key?: string;
    // id: number;
    // email?: string;
    // role:string;
    let roleVal = 2
    if (values.role == roles[3]) {
      roleVal = 3
    }else if (values.role == roles[1]) {
      roleVal = 1
    }
    const result = await serviceFun({
      user:{
        id,
        name:values.name,
        avatar:"",
        email:values.email,
        role: roleVal
      }
    });
    if (result) {
      setModalVisible(false);
      message.success(`${id === 0 ? '新增' : '更改'} 成功。`);
      resetHandler();
      setConfirmLoading(false);
    } else {
      setConfirmLoading(false);
      message.error(`${id === 0 ? '新增' : '更改'} 失败。`);
    }
  };

  const addHandler = () => {
    setModalVisible(true);
    setRecord(initUserInfo);
  };

  const resetHandler = () => {
    if (resetDataSource) {
      setResetDataSource(false)
    }else{
      setResetDataSource(true)
    }
    // dispatch({
    //   type: 'users/getRemote',
    //   payload: {
    //     page: users.meta.page,
    //     per_page: users.meta.per_page,
    //   },
    // });
  };

  const paginationHandler = (page: number, pageSize?: number) => {
    // dispatch({
    //   type: 'users/getRemote',
    //   payload: {
    //     page,
    //     per_page: pageSize ? pageSize : users.meta.per_page,
    //   },
    //});
  };

  const pageSizeHandler = (current: number, size: number) => {
    // dispatch({
    //   type: 'users/getRemote',
    //   payload: {
    //     page: current,
    //     per_page: size,
    //   },
    // });
  };

  return (
    <div className="list-table">
      <ProTable
        columns={columns}
        dataSource={dataSource}
        rowKey="id"
        
        //loading={userListLoading}
        search={{
          
        }}
        pagination={
          {
            pageSize: 3,
            hideOnSinglePage: true,
            total:dataSource.length
          }
        }
        options={{
          density: true,
          fullScreen: true,
          reload: () => {
            resetHandler();
          },
          setting: true,
        }}
        headerTitle="User List"
        toolBarRender={() => [
          <Button type="primary" onClick={addHandler}>
            新增
          </Button>,
          <Button onClick={resetHandler}>刷新</Button>,
        ]}
      />
      {/* <Pagination
        className="list-page"
        total={dataSource.length}
        onChange={paginationHandler}
        onShowSizeChange={pageSizeHandler}
        current={page}
        pageSize={pageSize}
        showSizeChanger
        showQuickJumper
        showTotal={total => `Total ${total} items`}
      /> */}
      <UserModal
        visible={modalVisible}
        closeHandler={closeHandler}
        record={record}
        onFinish={onFinish}
        confirmLoading={confirmLoading}
      ></UserModal>
    </div>
  );
};

// const mapStateToProps = ({
//   users,
//   loading,
// }: {
//   users: UserState;
//   loading: Loading;
// }) => {
//   return {
//     users,
//     userListLoading: loading.models.users,
//   };
// };

export default UserListPage;
