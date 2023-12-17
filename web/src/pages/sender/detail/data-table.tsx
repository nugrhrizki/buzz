import {
  ColumnDef,
  ColumnFiltersState,
  SortingState,
  VisibilityState,
  createSolidTable,
  getCoreRowModel,
  getExpandedRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
} from "@tanstack/solid-table";
import { Match, ParentProps, Switch, createSignal } from "solid-js";

import { generateAvatarUrl } from "@/lib/utils";

import { useAvatar } from "@/services/sender";

import { Contact } from "@/models/sender";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import DataTable from "@/components/ui/data-table";
import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table";

type ColumnDefiniton = ColumnDef<Contact>;

export function columns(token: string): ColumnDefiniton[] {
  return [
    DataTable.RowExpand(),
    {
      accessorKey: "avatar",
      header: (header) => <DataTable.ColumnHeader column={header.column} title="Avatar" />,
      cell: (cell) => {
        const avatar = useAvatar(
          () => token,
          () => cell.row.original.jid,
        );

        return (
          <Avatar>
            <AvatarImage
              src={
                avatar.data?.url ||
                generateAvatarUrl({
                  name: cell.row.original.jid,
                  style: "big-smile",
                  backgroundColors: ["#d6e6ff", "#d7f9f8", "#ffffea", "#fff0d4", "#fbe0e0", "#e5d4ef"],
                })
              }
            />
            <AvatarFallback>
              {cell.row.original.pushName
                .split(" ")
                .map((name) => name[0])
                .slice(0, 2)
                .join("")}
            </AvatarFallback>
          </Avatar>
        );
      },
      enableSorting: false,
    },
    {
      id: "name",
      accessorFn: (row) => {
        if (row.fullName === "") {
          return "~" + row.pushName;
        }
        return row.fullName;
      },
      header: (header) => <DataTable.ColumnHeader column={header.column} title="Nama" />,
      cell: (cell) => {
        if (cell.row.original.fullName === "") {
          return "~" + cell.row.original.pushName;
        }
        return cell.row.original.fullName;
      },
      enableHiding: false,
    },
    {
      accessorKey: "jid",
      header: (header) => <DataTable.ColumnHeader column={header.column} title="JID" />,
      enableSorting: false,
    },
    {
      accessorKey: "found",
      header: (header) => <DataTable.ColumnHeader column={header.column} title="Found" />,
      cell: (cell) => (
        <Switch fallback={<Badge variant="destructive">unknown</Badge>}>
          <Match when={cell.row.original.found === true}>
            <Badge variant="default" class="bg-green-600 hover:bg-green-500">
              on whatsapp
            </Badge>
          </Match>
          <Match when={cell.row.original.found === false}>
            <Badge variant="secondary">not on whatsapp</Badge>
          </Match>
        </Switch>
      ),
      enableSorting: false,
    },
  ];
}

export function createContactsTable(data: Contact[], columns: ColumnDef<Contact>[]) {
  const [sorting, setSorting] = createSignal<SortingState>([]);
  const [columnFilters, setColumnFilters] = createSignal<ColumnFiltersState>([]);
  const [columnVisibility, setColumnVisibility] = createSignal<VisibilityState>({
    jid: false,
  });
  const [rowSelection, setRowSelection] = createSignal({});

  return createSolidTable({
    data,
    columns,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getExpandedRowModel: getExpandedRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getRowCanExpand: () => true,
    onColumnVisibilityChange: setColumnVisibility,
    onRowSelectionChange: setRowSelection,
    state: {
      get sorting() {
        return sorting();
      },
      get columnFilters() {
        return columnFilters();
      },
      get columnVisibility() {
        return columnVisibility();
      },
      get rowSelection() {
        return rowSelection();
      },
    },
  });
}

interface TableDetailProps {
  contact: Contact;
}

export function TableDetail(props: ParentProps<TableDetailProps>) {
  return (
    <div class="p-4">
      <Card>
        <CardContent class="p-2 pt-2">
          <Table>
            <TableBody>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">Name</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>{props.contact.pushName || "-"}</TableCell>
              </TableRow>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">First Name</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>{props.contact.firstName || "-"}</TableCell>
              </TableRow>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">Full Name</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>{props.contact.fullName || "-"}</TableCell>
              </TableRow>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">Business Name</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>{props.contact.businessName || "-"}</TableCell>
              </TableRow>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">JID</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>{props.contact.jid || "-"}</TableCell>
              </TableRow>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">Found</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>
                  <Switch fallback={<Badge variant="destructive">unknown</Badge>}>
                    <Match when={props.contact.found === true}>
                      <Badge variant="default">on whatsapp</Badge>
                    </Match>
                    <Match when={props.contact.found === false}>
                      <Badge variant="secondary">not on whatsapp</Badge>
                    </Match>
                  </Switch>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
