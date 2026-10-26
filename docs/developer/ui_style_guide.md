<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->

# Style guide

This guide intends to help developers to implement features for ISDuBA that have a consistent look and behavior.

In general it is always important to check if elements work in the light *and* the dark theme. Furthermore, it is helpful to look for existing components that fulfills the rules of this guide and reuse those components. There are aleady some for SSVC and CVSS labels for example.

## Buttons

### Colors

#### Blue: Primary actions/Next step

Use blue color for buttons that trigger the most important actions on a page, e.g. search for a document or create a source or an aggregator.

If a button saves, updates, or deletes information in the backend prefer green or read one instead (see below).

If there is already a green button next to a button that would be a candidate for a blue button use a light button instead because the green button is more important and should not need to draw more attention to itself. Example: When you edit an SSVC there are two buttons: "Evaluate" and "Save". The latter is green because it saves the SSVC in the backend. The button "Evaluate" could be blue because it is an important action but since it is next to a green button it is a transparent one.

#### Green: Save/Update information

Whenever a user can save or update something like an SSVC or a comment offer a green button.

#### Red: Delete information permanently

When a button directly leads to the permanent deletion of a document etc. indicate this with a red button. When a button opens a confirmation dialog use a transparent button instead.

#### Transparent: Secondary controls

All buttons that do not belong to any of the groups above have to have a transparent background.

##### Transparent with blue outline

These are for the options in processes with multiple steps like the SSVC calculator.

### Borders

All buttons should have a border and it should be rounded. The rounded border is important to distinguish between buttons and labels. Only if there is a button next to it with no space between the border on this side should not be rounded.

### Spacing between buttons

If buttons belong to a group and work as radio buttons there has to be no space between them, e.g. the buttons to set the type at the search page. Otherwise there has to be a gap between them.

### Sizing

#### Height

Use two variants:

- Small: 28px
  - Border top: 1px
  - Padding top: 4px
  - Content: 18px
  - Padding bottom: 4px
  - Border bottom: 1px

- Big: 38px
  - Border top: 1px
  - Padding top: 8px
  - Content: 20px
  - Padding bottom: 8px
  - Border bottom: 1px

### Misc

Always display a hand cursor (CSS: `cursor: pointer`) when the user hovers over an active button. When buttons are disabled show the default one.

## Labels

**Note: CVSS and SSVC contain crucial information about a CSAF document. Therefore they should stand out from the rest of the document.**

### Colors

#### TLPs

Text and background color of TLP labels are defined by FIRST: https://www.first.org/tlp/.

#### SSVC

Here the background colors are given in the decision tree. An appropriate text color has to be chosen so that the text is readable well enough.

#### CVSS

Use the following colors for the background of CVSS labels and choose a well readable text color:

| Label    | Color   |
| -------- | ------- |
| None     | #53aa33 |
| Low      | #ffcb0d |
| Medium   | #f9a009 |
| High     | #df3d03 |
| Critical | #cc0500 |

### Borders

Labels with borders must not have rounded corners.

### Sizing

Labels for CVSS, SSVC, and TLP should be always 24px high and have a minimum width of 30px.

### Fonts

Use a semi-bold font for CVSS and SSVC and a normal font weight for TLP labels. While TLP labels should have a font size of 12 pt (https://www.first.org/tlp/) use a smaller semi-bold font for CVSS and CVSS.

### Format of dates

Prefer the [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format for dates.

- YYYY-MM-DD for the release dates (initial, current etc.) of an advisory **in the table with search results**
- YYYY-MM-DD hh:mm:ssZ for the release dates (initial, current etc.) of an advisory **in advisory view**
- Aug 25, 2026, 8:08:48 AM

#### Statistics

The format depends on the interval for which data is shown in the graph:

- MM-DD | hh:mm: Interval of maximal a week
- YYYY-MM-DD: Interval between a week and a year
- YYYY-Www: Interval of a year or longer

## Text

### Colors

Use the Tailwind classes `text-gray-900` or `text-gray-500` in the light and `text-white` or `text-gray-400` in the dark mode for text in cards and paragraphs. Exceptions may be texts inside badges.

For feeds in the aggregator that are not configured use `text-amber-600` in light and dark mode.

## Links

Internal links have to be blue, while external links have to be in a neutral color (white/light gray in dark mode, black/dark gray in light mode).

Links that are part of an advisory shall be presented as `<a>` elements and not normal text elements so users can open them easily.

## Icons

The same icon should be used for the same action everytime. Elements that always have to contain an icon are the top-level links in the sidebar and blue and green buttons. When an icon is used inside a button or next to a link it should be on the left side.

Never use an icon for secondary controls. Exceptions are buttons that do not contain any text. The workflow state buttons in the advisory view also play a special role because we want to display the workflow state icons everywhere where these states appear so the users can learn what the icons mean and know about their meaning when they see the icon only, e.g. in the table with the search result.

### Common icons

These icons are used in multiple areas of the client so it is even more important to ensure that they are used consistently:

- [Plus](https://boxicons.com/icons/plus?s=regular&w=normal&p=basic) for buttons to create an item.
- [Save](https://boxicons.com/icons/save?s=regular&w=normal&p=basic) for buttons to save an item.
  - Exception: [Send](https://boxicons.com/icons/send?s=regular&w=normal&p=basic) for buttons that create new comments.
- (Red) [Trash](https://boxicons.com/icons/trash?s=regular&w=normal&p=basic) for buttons to delete an item.
- [Arrow out up right stroke square](https://boxicons.com/icons/arrow-out-up-right-square?s=regular&w=normal&p=basic) for elements that lead to an external website. Either inside a button or next to a link.
- [X](https://boxicons.com/icons/x?s=regular&w=normal&p=basic) for buttons to close dialogs or hints.
- [Git repo fork](https://boxicons.com/icons/git-repo-forked?s=regular&w=normal&p=basic) to indicate that an item is a souce.
- [Sitemap](https://boxicons.com/icons/sitemap?s=regular&w=normal&p=basic) to indicate that an item is an aggregator.
- [List ul square](https://boxicons.com/icons/list-ul-square?s=regular&w=normal&p=basic) for links that lead to the search page.
- [Arrow down a z](https://boxicons.com/icons/arrow-down-a-z?s=regular&w=normal&p=basic) and [Arrow up a z](https://boxicons.com/icons/arrow-up-a-z?s=regular&w=normal&p=basic) to indicate that a list/table column can be sorted and how it is sorted currently.

## Tables

### Colors

In the light mode use `text-gray-700` for table headers and `text-gray-500` for the text inside the table body. The rows should be striped with `bg-white` and `bg-gray-100`.

In the dark mode use `text-gray-400` for the table headers and `text-gray-400` for the table body. The rows should be striped with `bg-gray-800` and `bg-gray-700`.

### Alignment

All table headers should be aligned to the left side while it depends for the table body cells. When a body cell always contains only a single symbol (icon, digit, ...) center the content. But when the content can be longer align it at the left side.

### Overflow

If a table grows out of the window on the horizontal axis let the user scroll the table without scrolling the whole page on this axis.

## Accordions

### Colors

The text inside accordion headers should be black in the light mode and white in the dark mode no matter if the accordion is open or not.

### Indicator for open/closed state

Always use a chevron icon to indicate the header and put it on the left side. When it would be on the right side and the title is very short there would be a large gap but elements that belong together should be positioned near to each other.

## Definition lists

In ISDuBA definition lists (`<dl>`) appear where meta information about an item are displayed, e.g. a source or an aggregator.

### Colors

In both themes use `text-gray-500` for the `<dt>` element but for `<dd>` elements use `text-gray-100` in the light and `text-gray-900` in the dark mode.

To divide the entries from each other use borders with `divide-gray-200` in the light and `divide-gray-600` in the dark theme. Do not use a border around the whole definition list.

### Fonts

The text of the `<dt>` elements should have a normal weight while the one of the `<dd>` elements should be semi-bold.

### Spacing

Add a margin of 0.25 rem on the y-axis to the lines that divide the entries.

## Errors

### Prevention/validation

Always validate the content of input fields. Extend labels with explanations about the format where necessary and if possible prefill input fields with reasonable values to prevent errors. In cases where the content is invalid it may make sense to automatically overwrite it (e.g. when the user enters an "a" in the page input of the pagination just overwrite it with the previous number). Otherwise display an error.

After cancelling a request to the backend it returns an `AbortError`. Do not show an error message for this error since it is expected that the request will not finish with a result.

### Placement

Put validation messages of the client right below the input field. When the user submits a form or presses a button to trigger an action and an error occurs place the error message near the trigger element.

## Animations

### Loading indicators

Loading indicators should always be spinners and they should fade in with a delay because some processes take only a short time and we want to avoid that there is only a short glimpse of a spinner. That would be annoying for users.

Some information in the UI is reloaded regularly to keep it up-to-date. To prevent surprises the UI shall show a blue line with 4 px height that is shrinking from right to left. Every time the line reaches a width of 0 the reload should be started.

### Other animations

When other elements appear in the UI after the user triggered an action they may slide in to prevent too much noise because of many elements changing their position.

When the content of a page is loaded you may use preview/skeleton elements (e.g. an empty disabled button instead of the buttons to choose a version in the advisory view) with a pulsating animation until they are ready. The reason is also to achieve a calm layout.

## Specialties

### Diff-box

When the user visits the search page (`/search`) they can click on a button with the label "Diff". With this click they open the diff-box and activate a mode where they can choose documents from the search result list. The mechanism works similar to many online shops where customers can select products to compare them. There the users can choose them from different pages too so they are familiar with this behavior.

### Version selection for diff

In the advisory view the client offers a button "Show changes" so the user can select two versions to create a diff of them. The currently opened and the previous version are pre-selected because the chance is high that the user wants to find out what the latest changes were. The diff in the lower part of the page shows the changes from the version marked with red color and the minus to the version marked with green color and the plus. The color alone is not enough as an indicator because some people cannot distinguish between red and green and other colors might be meaningless in this case.

When the user marked a version with a minus they shall not be able to mark the previous versions with a plus. Also, the last button shall not be enabled when the user has to mark a version with the minus. While the cursor is hovering over an enabled version button the color has to change slightly and the appropriate icon below the button has to
indicate with which icon the version will be marked. A circle has to surround the icon because otherwise a color-blind person could not distinguish e.g. between a version button marked with a minus and the same button but hovered by the cursor.

### Webview

The webview is a component that displays a general section about a document and specific information that is divided into several tabs. On a small screen only one tab can be presented to the user at a time while on wider screens multiple tabs should be stacked horizontally to use the available space efficiently.
