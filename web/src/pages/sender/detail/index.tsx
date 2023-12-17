import { useParams } from "@solidjs/router";
import { Table as TanstackTable } from "@tanstack/solid-table";
import { RiSystemRefreshLine } from "solid-icons/ri";
import { TbLoader } from "solid-icons/tb";
import { For, Show, createEffect, createSignal, onCleanup } from "solid-js";

import { connect, getContacts, getQR, statusRequest, useSenderDetail } from "@/services/sender";

import { Contact } from "@/models/sender";

import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion";
import { Button } from "@/components/ui/button";
import DataTable from "@/components/ui/data-table";
import { Input } from "@/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

import { QRCode, SenderInfo } from "@/components/sender/detail";

import { STATUS } from "./constant";
import { TableDetail, columns, createContactsTable } from "./data-table";

function SenderDetailPage() {
  const params = useParams();
  const sender = useSenderDetail(parseInt(params.id));
  const [scanInterval, setScanInterval] = createSignal<NodeJS.Timeout | null>(null);
  const [scanned, setScanned] = createSignal<boolean>(false);
  const [qrcode, setQrcode] = createSignal<string>("");
  const [status, setStatus] = createSignal<string>("");
  const [contacts, setContacts] = createSignal<Contact[]>([]);
  const [table, setTable] = createSignal<TanstackTable<Contact>>(
    createContactsTable(contacts(), columns(sender.data?.token || "")),
  );

  onCleanup(() => {
    if (scanInterval() !== null) {
      clearInterval(scanInterval()!);
      setScanInterval(null);
    }

    setScanned(false);
    setQrcode("");
    setStatus("");
  });

  createEffect(() => {
    if (status() === STATUS.CONNECTED) {
      getContacts(sender.data?.token || "").then((data) => {
        if (data === null) {
          return;
        }
        const contacts: Contact[] = [];
        for (const [key, value] of Object.entries(data)) {
          contacts.push({
            jid: key,
            found: value.Found,
            firstName: value.FirstName,
            fullName: value.FullName,
            pushName: value.PushName,
            businessName: value.BusinessName,
          });
        }

        setContacts(contacts);
      });
    }
  });

  createEffect(() => {
    setTable(createContactsTable(contacts(), columns(sender.data?.token || "")));
  });

  createEffect(() => {
    statusRequest(sender.data?.token || "").then((status) => {
      if (status === null) {
        setStatus(STATUS.NOT_CONNECTED);
        return;
      }

      if (status.success == true) {
        if (status.data.logged_in === false) {
          if (status.data.connected === true) {
            showQR();
          } else {
            setStatus(STATUS.NOT_CONNECTED);
            connect(sender.data?.token || "").then(() => {
              setStatus(STATUS.CONNECTED);
            });
          }
        } else {
          if (status.data.connected === false) {
            connect(sender.data?.token || "").then(() => {
              setStatus(STATUS.CONNECTED);
            });
          }
          setStatus(STATUS.CONNECTED);
          setScanned(true);
        }
      } else if (status.success == false) {
        if (status.error == "no session") {
          connect(sender.data?.token || "").then((data) => {
            if (data.success === true) {
              showQR();
            } else {
              setStatus(STATUS.COULD_NOT_CONNECT);
            }
          });
        } else if (status.error == "Unauthorized") {
          setStatus(STATUS.BAD_AUTHENTICATION);
        }
      } else {
        setStatus(STATUS.BAD_AUTHENTICATION);
      }
      return;
    });
  });

  function checkStatus() {
    statusRequest(sender.data?.token || "").then((status) => {
      if (status === null) {
        setStatus(STATUS.NOT_CONNECTED);
        return;
      }

      if (status.success == true) {
        if (status.data.logged_in === true) {
          setScanned(true);
          setStatus(STATUS.CONNECTED);
          if (scanInterval() !== null) {
            clearInterval(scanInterval()!);
            setScanInterval(null);
          }
        }
      } else {
        if (scanInterval() !== null) {
          clearInterval(scanInterval()!);
          setScanInterval(null);
        }
      }
    });
  }

  async function wait(time: number) {
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve(true);
      }, time);
    });
  }

  async function showQR() {
    if (scanInterval() !== null) {
      clearInterval(scanInterval()!);
      setScanInterval(null);
    }
    setScanInterval(setInterval(checkStatus, 1000));
    setStatus(STATUS.SHOW_QR);
    while (!scanned()) {
      setStatus(STATUS.NOT_SCANNED);
      var data = await getQR(sender.data?.token || "");
      if (data === null) {
        setStatus(STATUS.NOT_CONNECTED);
        return;
      }

      if (data.success == true) {
        setQrcode(data.qrcode);
        if (data.qrcode != "") {
          await wait(15 * 1000);
        }
      } else {
        setScanned(true);
        if (scanInterval() !== null) {
          clearInterval(scanInterval()!);
          setScanInterval(null);
        }

        setStatus(STATUS.TIMEOUT);
      }
    }
  }

  function senderData() {
    return sender.data;
  }

  return (
    <div class="space-y-4 p-8 pt-6">
      <div class="flex items-center justify-between space-y-2">
        <h2 class="text-3xl font-bold tracking-tight">
          Sender Detail
          <Show when={sender.isFetching}>
            <TbLoader class="inline-flex ml-2 h-4 w-4 animate-spin" />
          </Show>
        </h2>
        <Button variant="outline" onClick={() => sender.refetch()}>
          <RiSystemRefreshLine
            class="w-5 h-5"
            classList={{
              "animate-spin": sender.isFetching,
            }}
          />
        </Button>
      </div>
      <div class="grid grid-cols-12 gap-4">
        <div class="col-span-7 row-span-2">
          <Tabs defaultValue="contacts" class="space-y-4">
            <TabsList>
              <TabsTrigger value="contacts">Contacts</TabsTrigger>
              <TabsTrigger value="send-message">Send Message</TabsTrigger>
            </TabsList>
            <TabsContent value="contacts" class="space-y-4">
              <div class="flex gap-x-4">
                <Input
                  placeholder="Filter name..."
                  value={(table().getColumn("name")?.getFilterValue() as string) ?? ""}
                  onInput={(event) => {
                    table().getColumn("name")?.setFilterValue(event.target.value);
                  }}
                  class="max-w-sm"
                />
                <div class="ml-auto flex gap-x-4">
                  <DataTable.ColumnVisibility table={table()} />
                </div>
              </div>
              <DataTable.Root>
                <DataTable.Table table={table()}>{(row) => <TableDetail contact={row.original} />}</DataTable.Table>
                <DataTable.Pagination table={table()} />
              </DataTable.Root>
            </TabsContent>
            <TabsContent value="send-message" class="space-y-4">
              <Accordion multiple={false} collapsible class="w-full">
                <For
                  each={[
                    {
                      id: "send-text",
                      name: "Send Text",
                      content: <div>Send Text</div>,
                    },
                    {
                      id: "send-document",
                      name: "Send Document",
                      content: <div>Send Document</div>,
                    },
                    {
                      id: "send-audio",
                      name: "Send Audio",
                      content: <div>Send Audio</div>,
                    },
                    {
                      id: "send-image",
                      name: "Send Image",
                      content: <div>Send Image</div>,
                    },
                    {
                      id: "send-sticker",
                      name: "Send Sticker",
                      content: <div>Send Sticker</div>,
                    },
                    {
                      id: "send-video",
                      name: "Send Video",
                      content: <div>Send Video</div>,
                    },
                    {
                      id: "send-contact",
                      name: "Send Contact",
                      content: <div>Send Contact</div>,
                    },
                    {
                      id: "send-location",
                      name: "Send Location",
                      content: <div>Send Location</div>,
                    },
                    {
                      id: "send-button",
                      name: "Send Button",
                      content: <div>Send Button</div>,
                    },
                    {
                      id: "send-list",
                      name: "Send List",
                      content: <div>Send List</div>,
                    },
                  ]}>
                  {(item) => (
                    <AccordionItem value={item.id}>
                      <AccordionTrigger>{item.name}</AccordionTrigger>
                      <AccordionContent>{item.content}</AccordionContent>
                    </AccordionItem>
                  )}
                </For>
              </Accordion>
            </TabsContent>
          </Tabs>
        </div>
        <div class="col-span-5">
          <SenderInfo sender={senderData()} status={status} />
        </div>
        <div class="col-span-5">
          <Show when={qrcode() !== ""}>
            <QRCode qrcode={qrcode()} />
          </Show>
        </div>
      </div>
    </div>
  );
}

export default SenderDetailPage;
