# Copyright (c) 2021 Cisco Systems, Inc. and its affiliates
# All rights reserved
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
"""fledge role composer."""

import logging
from queue import Queue
from types import TracebackType
from typing import Optional, Type

logger = logging.getLogger(__name__)


class Composer(object):
    """Composer enables composition of tasklets."""

    #
    def __init__(self) -> None:
        """Initialize the class."""
        # maintain tasklet chains
        self.chain = dict()
        self.reverse_chain = dict()

    def __enter__(self):
        """Enter custom context."""
        ComposerContext.set_composer(self)
        return self

    def __exit__(
        self, exc_type: Optional[Type[BaseException]],
        exc_value: Optional[BaseException],
        exc_traceback: Optional[TracebackType]
    ) -> bool:
        """Exit custom context."""
        ComposerContext.reset_composer()

    def _get_tasklets_in_loop(self, start, end):
        tasklets_in_loop = set()

        # traverse tasklets and execute them
        q = Queue()
        q.put(start)
        visited = set()
        visited.add(start)
        while not q.empty():
            tasklet = q.get()

            tasklets_in_loop.add(tasklet)

            if tasklet is end:
                break

            for child in self.chain[tasklet]:
                if child not in visited:
                    visited.add(child)
                    q.put(child)

        return tasklets_in_loop

    def run(self):
        """Execute tasklets in the chain."""
        # choose one tasklet
        tasklet = next(iter(self.chain))
        root = tasklet.get_root()

        # traverse tasklets and execute them
        q = Queue()
        q.put(root)
        visited = set()
        visited.add(root)
        while not q.empty():
            tasklet = q.get()

            # execute tasklet
            tasklet.do()

            if tasklet.is_last_in_loop() and not tasklet.is_loop_done():
                start, end = tasklet.loop_starter, tasklet
                tasklets_in_loop = self._get_tasklets_in_loop(start, end)

                visited = visited - tasklets_in_loop
                tasklet = tasklet.loop_starter

            for child in self.chain[tasklet]:
                if child not in visited:
                    visited.add(child)
                    q.put(child)


class ComposerContext(object):
    """ComposerContext maintains a context of composer."""

    _context_composer: Optional[Composer] = None

    @classmethod
    def get_composer(cls) -> Optional[Composer]:
        """get_composer returns a composer."""
        return cls._context_composer

    @classmethod
    def set_composer(cls, composer: Composer) -> None:
        """set_composer set a new composer."""
        cls._context_composer = composer

    @classmethod
    def reset_composer(cls) -> None:
        """reset_composer set a composer None."""
        cls._context_composer = None