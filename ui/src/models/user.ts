import { history } from 'umi';
import {
  userLogin
} from '@/services/user';
import type { Effect, Reducer, ConnectProps } from 'umi';
import handleRedirect from '@/utils/handleRedirect';
import handleGetRootSubmenuKeys from '@/utils/handleGetRootSubmenuKeys';
import handleGetEachDatumFromNestedDataByKey from '@/utils/handleGetEachDatumFromNestedDataByKey';
import handleGetIndexValidMenuItemByPath from '@/utils/handleGetIndexValidMenuItemByPath';

import { userLogoutData,adminAuthList,makeUserMenu,internalMenu,adminMenu, initUserInfo ,makeUserAuthority, powerAuthList, internalAuthList, powerMenu, avatarVal} from './data';
/**
 * 全局用户数据
 * @description isLogin 是否登录过
 * @description data 用户信息
 * @description menu 菜单数据
 * @description authority 权限数组
 * @description loginBtnLoading 按钮loading状态
 * @description rootSubmenuKeys 子菜单的父级菜单key
 * @description layoutWrapperLoading 布局外层容器loading状态
 * @description indexAllMenuItemById 通过id索引索引菜单项的映射
 * @description indexAllMenuItemByPath 通过path索引所有菜单项的映射
 * @description indexValidMenuItemByPath 通过path索引有效菜单项的映射
 */
type UserModelState = {
  isLogin: boolean;
  data: API.UserInfo;
  menu: API.MenuData;
  authority: string[];
  loginBtnLoading: boolean;
  rootSubmenuKeys: React.Key[];
  layoutWrapperLoading: boolean;
  indexAllMenuItemById: IndexAllMenuItemByKey<'id'>;
  indexValidMenuItemByPath: IndexValidMenuItemByPath;
  indexAllMenuItemByPath: IndexAllMenuItemByKey<'path'>;
}

export type UserConnectedProps = {
  user: UserModelState;
} & ConnectProps;

type UserModelType = {
  namespace: 'user';
  state: UserModelState;
  effects: {
    login: Effect;
    logout: Effect;
    resetLoginStatus: Effect;
    
  };
  reducers: {
    save: Reducer<UserModelState>;
  };
}


const UserModel: UserModelType = {
  namespace: 'user',
  state: {
    data: initUserInfo,
    authority: [],
    isLogin: false,
    rootSubmenuKeys: [],
    loginBtnLoading: false,
    layoutWrapperLoading: false,
    menu: [
      {
        id: 1,
        key: '1',
        path: '/',
        label: '首页',
        redirect: '',
      },
    ],
    indexAllMenuItemById: {
      1: {
        id: 1,
        key: '1',
        path: '/',
        label: '首页',
        redirect: '',
      },
    },
    indexAllMenuItemByPath: {
      '/': {
        id: 1,
        key: '1',
        path: '/',
        label: '首页',
        redirect: '',
      },
    },
    indexValidMenuItemByPath: {
      '/': {
        id: 1,
        key: '1',
        path: '/',
        label: '首页',
        redirect: '',
      },
    },
  },
  effects: {
    //登录
    *login({ payload }, { call, put }) {

      yield put({
        type: 'save',
        payload: {
          loginBtnLoading: true,
        },
      });

      const res: API.LoginResponse = yield call(userLogin, payload);

      //登录成功之后设置token并获取用户信息等数据
      if (!res.code) {
        localStorage.setItem('Authorization', res.data.token);
        // let userInfoRes: API.UserInfoResponse = {
        //   data: initUserInfo,
        //   code: 0,
        //   message: '',
        // };
  
        let userAuthorityRes: API.UserAuthorityResponse = {
          data: {
            authority: [],
          },
          code: 0,
          message: '',
        };
  
        let menuRes: API.MenuDataResponse = {
          data: [],
          code: 0,
          message: '',
        };

        let userInfoRes = res.data.user
        userInfoRes.avatar = avatarVal
        console.log("===mode user ", adminAuthList,adminMenu)
        if (userInfoRes.role == 1) {
          userAuthorityRes = makeUserAuthority(adminAuthList)
          menuRes = makeUserMenu(adminMenu)
        }else if (userInfoRes.role == 2 ) {
          userAuthorityRes = makeUserAuthority(powerAuthList)
          menuRes = makeUserMenu(powerMenu)
        }else if (userInfoRes.role == 3 ){
          userAuthorityRes = makeUserAuthority(internalAuthList)
          menuRes = makeUserMenu(internalMenu)
        }
        
        

        const indexAllMenuItemByPath = handleGetEachDatumFromNestedDataByKey(menuRes.data, 'path');
      const indexValidMenuItemByPath = handleGetIndexValidMenuItemByPath(menuRes.data);

      //在登录完获取菜单数据之后做是否需要重定向的操作
      yield call(
        handleRedirect,
        window.location.pathname === '/',
        indexAllMenuItemByPath,
        indexValidMenuItemByPath,
      );

      yield put({
        type: 'save',
        payload: {
          isLogin: true,
          menu: menuRes.data,
          data: userInfoRes,
          loginBtnLoading: false,
          indexAllMenuItemByPath,
          indexValidMenuItemByPath,
          layoutWrapperLoading: false,
          authority: userAuthorityRes.data.authority,
          rootSubmenuKeys: handleGetRootSubmenuKeys(menuRes.data),
          indexAllMenuItemById: handleGetEachDatumFromNestedDataByKey(menuRes.data, 'id'),
        },
      });

      //为保证所有语句都return, 因此这里加一句这个
      return true;
      } else {
        yield put({
          type: 'save',
          payload: {
            loginBtnLoading: false,
          },
        });
      }
    },
   
    //登出
    *logout({ payload }, { call, put }) {
      //const res: API.LogoutResponse = yield call(userLogout, payload);
      const res = userLogoutData
      if (!res.code) {
        yield put({
          type: 'resetLoginStatus',
        });
      }
    },
    //重置登录状态
    *resetLoginStatus(_, { put }) {
      localStorage.removeItem('Authorization');
      
      yield put({
        type: 'save',
        payload: {
          isLogin: false,
          loginBtnLoading: false,
        },
      });
      
      //当前页面不是登录页, 且没有redirect参数的时候再设置redirect参数
      //而且这个redirect参数要包含url pathname和查询字符串参数, 即pathname及其后面的所有字符
      if (window.location.pathname !== '/login') {
        if (!/redirect=/.test(window.location.search)) {
          const redirectValue = `${window.location.pathname}${window.location.search}`;
          history.push(`/?redirect=${encodeURIComponent(redirectValue)}`);
        }
      }
      //history.push("/")
    },
  },
  reducers: {
    save(state, action) {
      return {
        ...state,
        ...action.payload,
      };
    },
  },
};

export default UserModel;
