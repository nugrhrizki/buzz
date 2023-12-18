import { createMutation, createQuery } from "@tanstack/solid-query";

import { request } from "@/lib/request";

import { AuthForm } from "@/pages/login/validations/auth";

import { User } from "@/models/user";

type Response<T> = {
  status: string;
  title: string;
  message: string;
  data: T;
};

const queryUser = () => request.get<Response<User>>("/api/v1/auth/identify");
const loginUser = (data: AuthForm) => request.post<Response<User>>("/api/v1/auth/login", data);

export function useUser() {
  return createQuery(() => ({
    queryKey: ["user"],
    queryFn: async () => {
      const response = await queryUser();
      return response.data;
    },
    retry: 0,
  }));
}

export function useLoginUser() {
  return createMutation(() => ({
    queryKey: ["user"],
    mutationKey: ["login"],
    mutationFn: async (data: AuthForm) => {
      const response = await loginUser(data);
      return response.data;
    },
    retry: 0,
  }));
}
