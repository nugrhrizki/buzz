import { createMutation, createQuery } from "@tanstack/solid-query";

import { request } from "@/lib/request";

import { CreateUser, User } from "@/models/user";

const queryUser = () => request.get<Array<User>>("/api/v1/user/get-all");
const createQueryUser = (user: CreateUser) => request.post<User>("/api/v1/user/create", user);
const updateQueryUser = (id: number, user: CreateUser) => request.put<User>(`/api/v1/user/update/${id}`, user);
const deleteQueryUser = (id: number) => request.delete<User>(`/api/v1/user/delete/${id}`);

export function useUsers() {
  return createQuery(() => ({
    queryKey: ["users"],
    queryFn: async () => {
      const response = await queryUser();
      return response.data;
    },
  }));
}

export function createUserMutation() {
  return createMutation(() => ({
    queryKey: ["users"],
    mutationKey: ["users"],
    mutationFn: async (users: CreateUser) => {
      const response = await createQueryUser(users);
      return response.data;
    },
  }));
}

export function updateUserMutation() {
  return createMutation(() => ({
    queryKey: ["users"],
    mutationKey: ["users"],
    mutationFn: async (props: { id: number; user: CreateUser }) => {
      const response = await updateQueryUser(props.id, props.user);
      return response.data;
    },
  }));
}

export function deleteUserMutation() {
  return createMutation(() => ({
    queryKey: ["users"],
    mutationKey: ["users"],
    mutationFn: async (id: number) => {
      const response = await deleteQueryUser(id);
      return response.data;
    },
  }));
}
